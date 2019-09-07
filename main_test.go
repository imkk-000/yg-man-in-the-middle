package main_test

import (
	"testing"
	. "yulgang"

	"github.com/stretchr/testify/assert"
)

func TestBuilderNewPacketInputCode1DataBytesShouldReturnNewPacket(t *testing.T) {
	expectedNewPacket := []byte{0x00, 0x01, 0x00, 0x03, 0x99, 0x98, 0x97}

	actualNewPacket := BuilderNewPacket(0x0001, []byte{0x99, 0x98, 0x97})

	assert.Equal(t, expectedNewPacket, actualNewPacket)
}

func TestInjectDataInputEmptyShouldReturnEmptyBytes(t *testing.T) {
	actualLen, actualData := InjectData(0, []byte{})

	assert.Zero(t, actualLen)
	assert.Empty(t, actualData)
}
