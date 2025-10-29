// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

var (
	withCreateDate    = []byte(`{"PolicyName":"readwrite","Policy":{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["admin:*"]},{"Effect":"Allow","Action":["s3:*"],"Resource":["arn:aws:s3:::*"]}]},"CreateDate":"2020-03-15T10:10:10Z","UpdateDate":"2021-03-15T10:10:10Z"}`)
	withoutCreateDate = []byte(`{"PolicyName":"readwrite","Policy":{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["admin:*"]},{"Effect":"Allow","Action":["s3:*"],"Resource":["arn:aws:s3:::*"]}]}}`)
)

func TestPolicyInfo(t *testing.T) {
	testCases := []struct {
		pi          *PolicyInfo
		expectedBuf []byte
	}{
		{
			&PolicyInfo{
				PolicyName: "readwrite",
				Policy:     []byte(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["admin:*"]},{"Effect":"Allow","Action":["s3:*"],"Resource":["arn:aws:s3:::*"]}]}`),
				CreateDate: time.Date(2020, time.March, 15, 10, 10, 10, 0, time.UTC),
				UpdateDate: time.Date(2021, time.March, 15, 10, 10, 10, 0, time.UTC),
			},
			withCreateDate,
		},
		{
			&PolicyInfo{
				PolicyName: "readwrite",
				Policy:     []byte(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["admin:*"]},{"Effect":"Allow","Action":["s3:*"],"Resource":["arn:aws:s3:::*"]}]}`),
			},
			withoutCreateDate,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run("", func(t *testing.T) {
			buf, err := json.Marshal(testCase.pi)
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(buf, testCase.expectedBuf) {
				t.Errorf("expected %s, got %s", string(testCase.expectedBuf), string(buf))
			}
		})
	}
}
