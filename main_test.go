package main_test

import (
	"testing"
	. "yulgang"
	"yulgang/model"

	"github.com/stretchr/testify/assert"
)

func TestBuilderNewPacketInputCode1DataBytesShouldReturnNewPacket(t *testing.T) {
	expectedNewPacket := []byte{0x00, 0x01, 0x03, 0x00, 0x99, 0x98, 0x97}

	actualNewPacket := BuilderNewPacket(0x0100, []byte{0x99, 0x98, 0x97})

	assert.Equal(t, expectedNewPacket, actualNewPacket)
}

func TestGetData8064Input8064BytesDataShouldReturnNew8064BytesData(t *testing.T) {
	expectedData, expectedConfig, expectedUser := []byte{0x09, 0x00, 0x31, 0x32, 0x37, 0x2E, 0x30, 0x2E, 0x30, 0x2E, 0x31, 0x58, 0x04, 0x04, 0x00, 0x66, 0x61, 0x6B, 0x65}, model.IpConfig{IP: "234.1.234.56", Port: 16000}, "fake"
	inputData := []byte{0x0C, 0x00, 0x32, 0x33, 0x34, 0x2E, 0x31, 0x2E, 0x32, 0x33, 0x34, 0x2E, 0x35, 0x36, 0x80, 0x3E, 0x04, 0x00, 0x66, 0x61, 0x6B, 0x65}

	actualData, actualConfig, actualUser := GetData8064(inputData, "127.0.0.1", 1112)

	assert.Equal(t, expectedData, actualData)
	assert.Equal(t, expectedConfig, actualConfig)
	assert.Equal(t, expectedUser, actualUser)
}

func TestInjectDataInputEmptyShouldReturnEmptyBytes(t *testing.T) {
	actualLen, actualData := InjectData(0, []byte{})

	assert.Zero(t, actualLen)
	assert.Empty(t, actualData)
}
