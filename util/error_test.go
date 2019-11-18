package util

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestEmptyCollectorHasNoErrors(t *testing.T) {
	collector := new(ErrorCollector)

	assert.Check(t, !collector.HasErrors())
}

func TestCollectorWithErrorHasErrors(t *testing.T) {
	collector := new(ErrorCollector)
	collector.Collect(fmt.Errorf("An error"))
	assert.Check(t, collector.HasErrors())
}

func TestCollectorShouldFlattenNestedCollectors(t *testing.T) {
	collector := new(ErrorCollector)
	nested := new(ErrorCollector)

	collector.Collect(fmt.Errorf("An error"))
	nested.Collect(fmt.Errorf("A nested error"))

	collector.Collect(nested)

	assert.ErrorContains(t, collector, "Error 1: A nested error")
}

func TestCollectorShouldNotCollectNilError(t *testing.T) {
	collector := new(ErrorCollector)
	collector.Collect(nil)

	assert.Check(t, !collector.HasErrors())
}

func TestCollectorShouldNotCollectEmptyNestedCollector(t *testing.T) {
	collector := new(ErrorCollector)
	nested := new(ErrorCollector)

	collector.Collect(nested)

	assert.Check(t, !collector.HasErrors())

}
