// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package store

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type SampleRecord struct {
	Version *int64
}

func (rec SampleRecord) GetVersion(recordJSON string) (int64, error) {
	record := &SampleRecord{}
	err := json.Unmarshal([]byte(recordJSON), record)
	if err != nil {
		return -1, errors.Wrapf(err, "Error unmarshaling record")
	}
	if record.Version == nil {
		return -1, errors.New("Record version is not set")
	}
	return *record.Version, nil
}

func TestValidateApplierNoRecordKey(t *testing.T) {
	applier := &STMApplier{
		record:     SampleRecord{},
		recordJSON: "record",
	}

	err := applier.validateApplier()
	assert.Error(t, err, "Expected error when record key is not set in applier")
}

func TestValidateApplierNoRecordJSON(t *testing.T) {
	applier := &STMApplier{
		record:    SampleRecord{},
		recordKey: "key",
	}

	err := applier.validateApplier()
	assert.Error(t, err, "Expected error when record JSON is not set in applier")
}

func TestValidateApplierNoRecord(t *testing.T) {
	applier := &STMApplier{
		recordKey:  "key",
		recordJSON: "record",
	}

	err := applier.validateApplier()
	assert.Error(t, err, "Expected error when record is not set in applier")
}

func TestAddRecordWhenNoOlderRecordExists(t *testing.T) {
	getKey := "key"

	newRecord := generateRecordWithVersion(t, 1)

	mockSTM := &mockSTM{
		getFunc: func(key string) string {
			assert.Equal(t, getKey, key, "Unexpected key for Get")
			return ""
		},
		putFunc: func(key string, val string, opts ...clientv3.OpOption) {
			assert.Equal(t, getKey, key, "Unexpected key in Put")
			assert.Equal(t, val, newRecord, "Unexpected record in Put")
		},
	}

	applier := &STMApplier{
		record:     SampleRecord{},
		recordKey:  "key",
		recordJSON: newRecord,
	}

	err := applier.applyRecord(mockSTM)
	assert.NoError(t, err, "Unexpected error adding a record when no older record exists")
}

func TestAddRecordWhenOlderRecordWithLowerVersionExists(t *testing.T) {
	getKey := "key"

	existingRecord := generateRecordWithVersion(t, 1)
	newRecord := generateRecordWithVersion(t, 2)

	mockSTM := &mockSTM{
		getFunc: func(key string) string {
			assert.Equal(t, getKey, key, "Unexpected key for Get")
			return existingRecord
		},
		putFunc: func(key string, val string, opts ...clientv3.OpOption) {
			assert.Equal(t, getKey, key, "Unexpected key in Put")
			assert.Equal(t, val, newRecord, "Unexpected record in Put")
		},
	}

	applier := &STMApplier{
		record:     SampleRecord{},
		recordKey:  "key",
		recordJSON: newRecord,
	}

	err := applier.applyRecord(mockSTM)
	assert.NoError(t, err, "Unexpected error adding a record when older record with lower version exists")
}

func TestAddRecordWhenOlderRecordWithHigherVersionExists(t *testing.T) {
	getKey := "key"

	existingRecord := generateRecordWithVersion(t, 2)
	newRecord := generateRecordWithVersion(t, 1)

	mockSTM := &mockSTM{
		getFunc: func(key string) string {
			assert.Equal(t, getKey, key, "Unexpected key for Get")
			return existingRecord
		},
	}

	applier := &STMApplier{
		record:     SampleRecord{},
		recordKey:  "key",
		recordJSON: newRecord,
	}

	err := applier.applyRecord(mockSTM)
	assert.NoError(t, err, "Unexpected error adding a record when no older record with higher version exists")
}

func TestAddRecordWhenOlderRecordIsInvalid(t *testing.T) {
	getKey := "key"

	existingRecord := "invalidJSON"
	newRecord := generateRecordWithVersion(t, 1)

	mockSTM := &mockSTM{
		getFunc: func(key string) string {
			assert.Equal(t, getKey, key, "Unexpected key for Get")
			return existingRecord
		},
	}

	applier := &STMApplier{
		record:     SampleRecord{},
		recordKey:  "key",
		recordJSON: newRecord,
	}

	err := applier.applyRecord(mockSTM)
	assert.Error(t, err, "Expected an error while adding a record when older record is invalid")
}

func TestAddRecordWhenNewerRecordIsInvalid(t *testing.T) {
	getKey := "key"

	existingRecord := generateRecordWithVersion(t, 1)
	newRecord := "invalidJSON"

	mockSTM := &mockSTM{
		getFunc: func(key string) string {
			assert.Equal(t, getKey, key, "Unexpected key for Get")
			return existingRecord
		},
	}

	applier := &STMApplier{
		record:     SampleRecord{},
		recordKey:  "key",
		recordJSON: newRecord,
	}

	err := applier.applyRecord(mockSTM)
	assert.Error(t, err, "Expected an error while adding a record when new record is invalid")
}

func TestAddRecordWhenOlderRecordHasNoVersion(t *testing.T) {
	getKey := "key"

	existingRecord := generateRecordWithNoVersion(t)
	newRecord := generateRecordWithVersion(t, 1)

	mockSTM := &mockSTM{
		getFunc: func(key string) string {
			assert.Equal(t, getKey, key, "Unexpected key for Get")
			return existingRecord
		},
	}

	applier := &STMApplier{
		record:     SampleRecord{},
		recordKey:  "key",
		recordJSON: newRecord,
	}

	err := applier.applyRecord(mockSTM)
	assert.Error(t, err, "Expected an error while adding a record when older record has no version")
}

func TestAddRecordWhenNewerRecordHasNoVersion(t *testing.T) {
	getKey := "key"

	existingRecord := generateRecordWithVersion(t, 1)
	newRecord := generateRecordWithNoVersion(t)

	mockSTM := &mockSTM{
		getFunc: func(key string) string {
			assert.Equal(t, getKey, key, "Unexpected key for Get")
			return existingRecord
		},
	}

	applier := &STMApplier{
		record:     SampleRecord{},
		recordKey:  "key",
		recordJSON: newRecord,
	}

	err := applier.applyRecord(mockSTM)
	assert.Error(t, err, "Expected an error while adding a record when new record has no version")
}

func generateRecordWithVersion(t *testing.T, version int64) string {
	rec, err := json.Marshal(
		&SampleRecord{
			Version: aws.Int64(version),
		})
	assert.NoError(t, err, "Error generating a json string for record")
	return string(rec)
}

func generateRecordWithNoVersion(t *testing.T) string {
	rec, err := json.Marshal(&SampleRecord{})
	assert.NoError(t, err, "Error generating a json string for record")
	return string(rec)
}
