// Copyright 2023 The Khulnasoft Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package khulnasoft

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsedPaging(t *testing.T) {
	resp := newResponse(&http.Response{
		Header: http.Header{
			"Link": []string{
				strings.Join(
					[]string{
						`<https://try.khulnasoft.io/api/v1/repos/khulnasoft/go-sdk/issues/1/comments?page=3>; rel="next"`,
						`<https://try.khulnasoft.io/api/v1/repos/khulnasoft/go-sdk/issues/1/comments?page=4>; rel="last"`,
						`<https://try.khulnasoft.io/api/v1/repos/khulnasoft/go-sdk/issues/1/comments?page=1>; rel="first"`,
						`<https://try.khulnasoft.io/api/v1/repos/khulnasoft/go-sdk/issues/1/comments?page=1>; rel="prev"`,
					}, ",",
				),
			},
		},
	})

	assert.Equal(t, 1, resp.FirstPage)
	assert.Equal(t, 1, resp.PrevPage)
	assert.Equal(t, 3, resp.NextPage)
	assert.Equal(t, 4, resp.LastPage)
}
