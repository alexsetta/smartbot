package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	asset, err := NewAsset("btcbrl")
	assert.Nil(t, err)
	assert.NotNil(t, asset)

	outJson, err := asset.Find()
	assert.Nil(t, err)
	assert.NotNil(t, outJson)

	assert.Equal(t, "BTCBRL", outJson.Simbolo, "BTCBRL should equal BTCBRL")
	assert.NotZero(t, outJson.Preco, "Preco should not be zero")

	fmt.Println(PrettyJson(outJson))
}

func TestGetAll(t *testing.T) {
	asset, err := NewAsset("all")
	assert.Nil(t, err)
	assert.NotNil(t, asset)

	outJson, err := asset.GetAll()
	assert.Nil(t, err)
	assert.NotNil(t, outJson)

	fmt.Println(PrettyJson(outJson))
}
