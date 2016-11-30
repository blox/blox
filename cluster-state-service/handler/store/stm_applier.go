// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"github.com/blox/blox/cluster-state-service/handler/types"
	log "github.com/cihub/seelog"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/pkg/errors"
)

type STMApplier struct {
	record     types.Record
	recordKey  string
	recordJSON string
}

// applyVersionedRecord adds a new record to the store if the
// version number in the record is higher than the one that exists in the
// store
func (applier STMApplier) applyVersionedRecord(stm concurrency.STM) error {
	err := applier.validateApplier()
	if err != nil {
		return err
	}

	// Get existing record
	existingRecord := stm.Get(applier.recordKey)

	// A record already exists. Add new record only of the newer record has a higher version.
	if existingRecord != "" {
		existingRecordVersion, err := applier.record.GetVersion(existingRecord)
		if err != nil {
			return errors.Wrapf(err,
				"Error retrieving the version of the existing record in the STM applier")
		}
		newRecordVersion, err := applier.record.GetVersion(applier.recordJSON)
		if err != nil {
			return errors.Wrapf(err,
				"Error retrieving the version of the new record in the STM applier")
		}
		if existingRecordVersion >= newRecordVersion {
			log.Debugf("Not adding record for key %s with version %d as version %d already exists",
				applier.recordKey, newRecordVersion, existingRecordVersion)

			// Higher or equivalent version of the event has already been stored.
			return nil
		}
	}

	// New record has a higher version. Add it.
	stm.Put(applier.recordKey, applier.recordJSON)
	return nil
}

// applyVersionedRecord adds a new unversioned record to the store.
// An unversioned record is generated while bootstrapping and reconciliation
// workflows.
func (applier STMApplier) applyUnversionedRecord(stm concurrency.STM) error {
	err := applier.validateApplier()
	if err != nil {
		return err
	}

	// Get existing record
	existingRecord := stm.Get(applier.recordKey)

	// A record already exists. Don't do anything.
	if existingRecord != "" {
		return nil
	}

	// No record exists. Add the unversioned record.
	stm.Put(applier.recordKey, applier.recordJSON)
	return nil
}

func (applier STMApplier) validateApplier() error {
	if applier.record == nil {
		return errors.New("Record has to be initialized for the STM applier")
	}
	if applier.recordKey == "" {
		return errors.New("Record key cannot be empty for the STM applier")
	}
	if applier.recordJSON == "" {
		return errors.New("Record JSON cannot be ampty for the STM applier")
	}
	return nil
}
