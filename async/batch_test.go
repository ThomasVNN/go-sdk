/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBatch(t *testing.T) {
	assert := assert.New(t)

	workItems := 32

	items := make(chan interface{}, workItems)
	for x := 0; x < workItems; x++ {
		items <- "hello" + strconv.Itoa(x)
	}

	var processed int32
	action := func(_ context.Context, v interface{}) error {
		atomic.AddInt32(&processed, 1)
		return fmt.Errorf("this is only a test")
	}

	errors := make(chan error, workItems)
	NewBatch(
		items,
		action,
		OptBatchErrors(errors),
		OptBatchParallelism(4),
	).Process(context.Background())

	assert.Equal(workItems, processed)
	assert.Equal(workItems, len(errors))
}

func TestBatchPanic(t *testing.T) {
	assert := assert.New(t)

	workItems := 32

	items := make(chan interface{}, workItems)
	for x := 0; x < workItems; x++ {
		items <- "hello" + strconv.Itoa(x)
	}

	var processed int32
	action := func(_ context.Context, v interface{}) error {
		if result := atomic.AddInt32(&processed, 1); result == 1 {
			panic("this is only a test")
		}
		return nil
	}

	errors := make(chan error, workItems)
	NewBatch(items, action, OptBatchErrors(errors)).Process(context.Background())

	assert.Equal(workItems, processed)
	assert.Equal(1, len(errors))
}
