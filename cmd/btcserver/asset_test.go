package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	asset, err := NewAsset("btcbrl")
	assert.Nil(t, err)
	assert.NotNil(t, asset)

	err = asset.Find()
	assert.Nil(t, err)
	assert.NotNil(t, asset.data)

	assert.Equal(t, "BTCBRL", asset.data.Simbolo, "BTCBRL should equal BTCBRL")
	assert.NotZero(t, asset.data.Preco, "Preco should not be zero")

	//fmt.Println(PrettyJson(asset.data))
}

func TestGetAll(t *testing.T) {
	asset, err := NewAsset("all")
	assert.Nil(t, err)
	assert.NotNil(t, asset)

	outJson, err := asset.GetAll()

	assert.Nil(t, err)
	assert.NotNil(t, outJson)
	assert.NotZero(t, len(outJson), "outJson should not be zero")
	assert.NotZero(t, outJson[0].Preco, "Preco should not be zero")

	//fmt.Println(PrettyJson(outJson))
}
