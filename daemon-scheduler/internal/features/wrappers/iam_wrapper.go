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

package wrappers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

type IAMWrapper struct {
	client *iam.IAM
}

func NewIAMWrapper() IAMWrapper {
	awsSession, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	return IAMWrapper{
		client: iam.New(awsSession),
	}
}

func (iamWrapper IAMWrapper) GetInstanceProfile(instanceProfileName *string) error {
	in := iam.GetInstanceProfileInput{
		InstanceProfileName: instanceProfileName,
	}

	_, err := iamWrapper.client.GetInstanceProfile(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to get instance profile '%v'. ", instanceProfileName)
	}

	return nil
}

func (iamWrapper IAMWrapper) GetRole(roleName *string) error {
	in := iam.GetRoleInput{
		RoleName: roleName,
	}

	_, err := iamWrapper.client.GetRole(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to get role '%v'. ", roleName)
	}

	return nil
}

func (iamWrapper IAMWrapper) DeleteInstanceProfile(instanceProfile *string) error {
	in := iam.DeleteInstanceProfileInput{
		InstanceProfileName: instanceProfile,
	}

	_, err := iamWrapper.client.DeleteInstanceProfile(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to delete instance profile '%v'. ", instanceProfile)
	}

	return nil
}

func (iamWrapper IAMWrapper) DeleteRole(roleName *string) error {
	in := iam.DeleteRoleInput{
		RoleName: roleName,
	}

	_, err := iamWrapper.client.DeleteRole(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to delete role '%v'. ", roleName)
	}

	return nil
}

func (iamWrapper IAMWrapper) CreateInstanceProfile(instanceProfileName *string) error {
	in := iam.CreateInstanceProfileInput{
		InstanceProfileName: instanceProfileName,
	}

	_, err := iamWrapper.client.CreateInstanceProfile(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to create instance profile '%v'. ", instanceProfileName)
	}

	return nil
}

func (iamWrapper IAMWrapper) RemoveRoleFromInstanceProfile(roleName *string, instanceProfileName *string) error {
	in := iam.RemoveRoleFromInstanceProfileInput{
		RoleName:            roleName,
		InstanceProfileName: instanceProfileName,
	}

	_, err := iamWrapper.client.RemoveRoleFromInstanceProfile(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to remove role '%v' from instance profile '%v'. ", roleName, instanceProfileName)
	}

	return nil
}

func (iamWrapper IAMWrapper) AddRoleToInstanceProfile(roleName *string, instanceProfileName *string) error {
	in := iam.AddRoleToInstanceProfileInput{
		RoleName:            roleName,
		InstanceProfileName: instanceProfileName,
	}

	_, err := iamWrapper.client.AddRoleToInstanceProfile(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to add role '%v' to instance profile '%v'. ", roleName, instanceProfileName)
	}

	return nil
}

func (iamWrapper IAMWrapper) CreateRole(roleName *string, assumeRolePolicyDocument *string) error {
	in := iam.CreateRoleInput{
		RoleName:                 roleName,
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
	}

	_, err := iamWrapper.client.CreateRole(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to create role '%v'. ", roleName)
	}

	return nil
}

func (iamWrapper IAMWrapper) AttachRolePolicy(policyARN *string, roleName *string) error {
	in := iam.AttachRolePolicyInput{
		PolicyArn: policyARN,
		RoleName:  roleName,
	}

	_, err := iamWrapper.client.AttachRolePolicy(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to attach policy to role '%v'. ", roleName)
	}

	return nil
}

func (iamWrapper IAMWrapper) DetachRolePolicy(policyARN *string, roleName *string) error {
	in := iam.DetachRolePolicyInput{
		PolicyArn: policyARN,
		RoleName:  roleName,
	}

	_, err := iamWrapper.client.DetachRolePolicy(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to detach policy to role '%v'. ", roleName)
	}

	return nil
}
