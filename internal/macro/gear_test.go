package macro

import (
	"testing"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/stretchr/testify/assert"
)

func Test_compact(t *testing.T) {
	assert.Equal(t, []oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
	}, compact([]oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
	}, 2))

	assert.Equal(t, []oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
	}, compact([]oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
	}, 2))

	assert.Equal(t, []oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
	}, compact([]oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
	}, 4))
}
