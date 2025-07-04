package salesforce

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jszwec/csvutil"
)

type IteratorJob interface {
	Next(ctx context.Context) bool
	Error(ctx context.Context) error
	Decode(any) error
}

type bulkJobQueryIterator struct {
	NumberOfRecords int    `json:"Sforce-Numberofrecords"`
	Locator         string `json:"Sforce-Locator"`
	auth            *authentication
	uri             string
	err             error
	reader          io.ReadCloser
	config          *configuration
}

func (sf *Salesforce) newBulkJobQueryIterator(
	ctx context.Context,
	bulkJobId string,
) (*bulkJobQueryIterator, error) {
	pollErr := sf.waitForJobResults(ctx, bulkJobId, queryJobType, (time.Second / 2))
	if pollErr != nil {
		return nil, pollErr
	}
	return &bulkJobQueryIterator{
		auth:   sf.auth,
		uri:    "/jobs/query/" + bulkJobId + "/results",
		config: sf.config,
	}, nil
}

func (it *bulkJobQueryIterator) Next(ctx context.Context) bool {
	if it.reader != nil {
		it.err = it.reader.Close()
		if it.Locator == "" {
			return false
		}
	}
	uri := it.uri
	if it.Locator != "" {
		uri += "/?locator=" + it.Locator
	}
	resp, err := doRequest(
		ctx,
		it.auth,
		it.config,
		requestPayload{
			method:   http.MethodGet,
			uri:      uri,
			content:  jsonType,
			compress: it.config.compressionHeaders,
		},
	)
	if err != nil {
		it.err = err
		return false
	}
	it.reader = resp.Body

	it.NumberOfRecords, _ = strconv.Atoi(resp.Header["Sforce-Numberofrecords"][0])
	if resp.Header["Sforce-Locator"][0] != "null" {
		it.Locator = resp.Header["Sforce-Locator"][0]
	} else {
		it.Locator = ""
	}

	return true
}

func (it *bulkJobQueryIterator) Decode(val any) error {
	dec, err := csvutil.NewDecoder(csv.NewReader(it.reader))
	if err != nil {
		return fmt.Errorf("NewDecoder: %w", err)
	}

	if err := dec.Decode(val); err != nil && err != io.EOF {
		return fmt.Errorf("Decode: %w", err)
	}
	return nil
}

func (it *bulkJobQueryIterator) Error(_ context.Context) error {
	return it.err
}
