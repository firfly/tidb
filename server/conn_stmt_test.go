// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/mysql"
	"github.com/pingcap/tidb/terror"
)

func (ts ConnTestSuite) TestParseStmtArgs(c *C) {
	type args struct {
		args        []interface{}
		boundParams [][]byte
		nullBitmap  []byte
		paramTypes  []byte
		paramValues []byte
	}
	tests := []struct {
		args   args
		err    error
		expect interface{}
	}{
		// Tests for date/datetime/timestamp
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{12, 0},
				[]byte{0x0b, 0xda, 0x07, 0x0a, 0x11, 0x13, 0x1b, 0x1e, 0x01, 0x00, 0x00, 0x00},
			},
			nil,
			"2010-10-17 19:27:30.000001",
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{10, 0},
				[]byte{0x04, 0xda, 0x07, 0x0a, 0x11},
			},
			nil,
			"2010-10-17",
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{7, 0},
				[]byte{0x0b, 0xda, 0x07, 0x0a, 0x11, 0x13, 0x1b, 0x1e, 0x01, 0x00, 0x00, 0x00},
			},
			nil,
			"2010-10-17 19:27:30.000001",
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{7, 0},
				[]byte{0x07, 0xda, 0x07, 0x0a, 0x11, 0x13, 0x1b, 0x1e},
			},
			nil,
			"2010-10-17 19:27:30",
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{7, 0},
				[]byte{0x00},
			},
			nil,
			"0",
		},
		// Tests for time
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{11, 0},
				[]byte{0x0c, 0x01, 0x78, 0x00, 0x00, 0x00, 0x13, 0x1b, 0x1e, 0x01, 0x00, 0x00, 0x00},
			},
			nil,
			"-120 19:27:30.000001",
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{11, 0},
				[]byte{0x08, 0x01, 0x78, 0x00, 0x00, 0x00, 0x13, 0x1b, 0x1e},
			},
			nil,
			"-120 19:27:30",
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{11, 0},
				[]byte{0x00},
			},
			nil,
			"0",
		},
		// For error test
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{7, 0},
				[]byte{10},
			},
			mysql.ErrMalformPacket,
			nil,
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{11, 0},
				[]byte{10},
			},
			mysql.ErrMalformPacket,
			nil,
		},
		{
			args{
				make([]interface{}, 1),
				[][]byte{nil},
				[]byte{0x0},
				[]byte{11, 0},
				[]byte{8, 2},
			},
			mysql.ErrMalformPacket,
			nil,
		},
	}
	for _, tt := range tests {
		err := parseStmtArgs(tt.args.args, tt.args.boundParams, tt.args.nullBitmap, tt.args.paramTypes, tt.args.paramValues)
		c.Assert(terror.ErrorEqual(err, tt.err), IsTrue)
		c.Assert(tt.args.args[0], Equals, tt.expect)
	}
}
