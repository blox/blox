package json

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type randomStruct struct {
	ID      string `json:"id"`
	Account string `json:"account"`
}

type MarshallerTestSuite struct {
	suite.Suite
	r     randomStruct
	rJSON string
}

func (testSuite *MarshallerTestSuite) SetupTest() {
	testSuite.r = randomStruct{
		ID:      "37e398c5-ef55-40c0-af6d-29694badc036",
		Account: "12345678910",
	}
	jsonR, err := json.Marshal(testSuite.r)
	assert.Nil(testSuite.T(), err, "Cannot set up marshaller test suite")

	testSuite.rJSON = string(jsonR)
}

func TestMarshallerTestSuite(t *testing.T) {
	suite.Run(t, new(MarshallerTestSuite))
}

func (testSuite *MarshallerTestSuite) TestUnmarshalEmptyString() {
	var r randomStruct
	err := UnmarshalJSON("", &r)
	assert.Error(testSuite.T(), err, "Expected an error when trying to unmarshal empty string")
}

func (testSuite *MarshallerTestSuite) TestUnmarshalIntoNilInterface() {
	err := UnmarshalJSON("test", nil)
	assert.Error(testSuite.T(), err, "Expected an error when trying to unmarshal into a nil interface")
}

func (testSuite *MarshallerTestSuite) TestUnmarshalInvalidJSON() {
	var r randomStruct
	err := UnmarshalJSON("test", &r)
	assert.Error(testSuite.T(), err, "Expected an error when trying to unmarshal invalid JSON")
}

func (testSuite *MarshallerTestSuite) TestUnmarshal() {
	var r randomStruct
	err := UnmarshalJSON(testSuite.rJSON, &r)
	assert.Nil(testSuite.T(), err, "Unexpected error when unmarshaling json")
	assert.Exactly(testSuite.T(), testSuite.r, r, "Unmarshal returned unexpected results")
}
