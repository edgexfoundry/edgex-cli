/*******************************************************************************
 * Copyright 2019 VMware Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package client

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var TestID = "testid"
var TestPathID = "id/"
var TestInvalidPathID = []rune{
	0x7f,
}
var TestInvalidPathName = []rune{
	0x7f,
}
var TestPathName = "name/"
var TestPort = "1234"
var TestVerboseTrue = true
var TestVerboseFalse = false

//TO:DO Test verbose
func TestGetAllItems(t *testing.T) {
	tests := []struct {
		name           string
		mockClient     *http.Client
		isVerbose      bool
		expectedResult []byte
		expectedError  bool
	}{
		{
			name:           "Successful call",
			mockClient:     mockClientGetAllSuccess(),
			isVerbose:      true,
			expectedResult: []byte(SUCCESSFUL_DELETE),
			expectedError:  false,
		},
		{
			name:           "Unsuccessful http Get",
			mockClient:     mockClientGetAllErr(),
			isVerbose:      true,
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client = test.mockClient
			actual, err := GetAllItems("blabla", "1312")

			if test.expectedError && err == nil {
				t.Error("Expected an error")
				return
			}

			if !test.expectedError && err != nil {
				t.Errorf("Unexpectedly encountered error: %s", err)
				return
			}

			if !reflect.DeepEqual(test.expectedResult, actual) {
				t.Errorf("Expected result does not match the observed.\nExpected: %v\nObserved: %v\n", test.expectedResult, actual)
				return
			}
		})
	}
}

func mockClientGetAllSuccess() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(SUCCESSFUL_DELETE))
	})
	httpclient, _ := testingHTTPClient(h)
	return httpclient
}

func mockClientGetAllErr() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
	})
	httpclient, _ := testingHTTPClientErr(h)
	return httpclient
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}
	return cli, s.Close
}

func testingHTTPClientErr(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := NewTestClient(func(req *http.Request) *http.Response {
		return nil
	})
	return cli, s.Close
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), errors.New("error test")
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func mockClientDeleteByIDSuccess() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(SUCCESSFUL_DELETE))
	})
	httpclient, _ := testingHTTPClient(h)
	return httpclient
}

func mockClientDeleteByIDErr() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	})
	httpclient, _ := testingHTTPClientErr(h)
	return httpclient
}

func mockClientDeleteByIDBufferErr() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	})
	httpclient, _ := testingHTTPClient(h)
	return httpclient
}

func TestDeleteItemByID(t *testing.T) {
	tests := []struct {
		name           string
		mockClient     *http.Client
		testID         string
		testPathID     string
		expectedResult []byte
		expectedError  bool
	}{
		{
			name:           "Successful call",
			mockClient:     mockClientDeleteByIDSuccess(),
			testID:         TestID,
			testPathID:     TestPathID,
			expectedResult: []byte(SUCCESSFUL_DELETE),
			expectedError:  false,
		},
		{
			name:           "Unsuccessful http new Request",
			mockClient:     mockClientDeleteByIDErr(),
			testID:         TestID,
			testPathID:     string(TestInvalidPathID),
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name:           "Unsuccessful http Do",
			mockClient:     mockClientDeleteByIDErr(),
			testID:         TestID,
			testPathID:     TestPathID,
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name:           "Unsuccessful ioutil ReadAll",
			mockClient:     mockClientDeleteByIDBufferErr(),
			testID:         TestID,
			testPathID:     TestPathID,
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client = test.mockClient
			actual, err := DeleteItemByID(test.testID, test.testPathID, TestPort)

			if test.expectedError && err == nil {
				t.Error("Expected an error")
				return
			}

			if !test.expectedError && err != nil {
				t.Errorf("Unexpectedly encountered error: %s", err)
				return
			}

			if !reflect.DeepEqual(test.expectedResult, actual) {
				t.Errorf("Expected result does not match the observed.\nExpected: %v\nObserved: %v\n", test.expectedResult, actual)
				return
			}
		})
	}
}

func mockClientDeleteByNameSuccess() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(SUCCESSFUL_DELETE))
	})
	httpclient, _ := testingHTTPClient(h)
	return httpclient
}

func mockClientDeleteByNameBufferErr() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	})
	httpclient, _ := testingHTTPClient(h)
	return httpclient
}

func mockClientDeleteByNameErr() *http.Client {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
	})
	httpclient, _ := testingHTTPClientErr(h)
	return httpclient
}

func TestDeleteItemByName(t *testing.T) {
	tests := []struct {
		name           string
		mockClient     *http.Client
		testID         string
		testPathName   string
		expectedResult []byte
		expectedError  bool
	}{
		{
			name:           "Successful call",
			mockClient:     mockClientDeleteByNameSuccess(),
			testID:         TestID,
			testPathName:   TestPathName,
			expectedResult: []byte(SUCCESSFUL_DELETE),
			expectedError:  false,
		},
		{
			name:           "Unsuccessful http new Request",
			mockClient:     mockClientDeleteByNameErr(),
			testID:         TestID,
			testPathName:   string(TestInvalidPathID),
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name:           "Unsuccessful http Do",
			mockClient:     mockClientDeleteByNameErr(),
			testID:         TestID,
			testPathName:   TestPathName,
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name:           "Unsuccessful ioutil ReadAll",
			mockClient:     mockClientDeleteByNameBufferErr(),
			testID:         TestID,
			testPathName:   TestPathName,
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client = test.mockClient
			actual, err := DeleteItemByName(test.testID, test.testPathName, TestPort)

			if test.expectedError && err == nil {
				t.Error("Expected an error")
				return
			}

			if !test.expectedError && err != nil {
				t.Errorf("Unexpectedly encountered error: %s", err)
				return
			}

			if !reflect.DeepEqual(test.expectedResult, actual) {
				t.Errorf("Expected result does not match the observed.\nExpected: %v\nObserved: %v\n", test.expectedResult, actual)
				return
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name           string
		mockClient     *http.Client
		testID         string
		testPathID     string
		testPathName   string
		expectedResult []byte
		expectedError  bool
	}{
		{
			name:         "Successful call",
			mockClient:   mockClientDeleteByNameSuccess(),
			testID:       TestID,
			testPathID:   TestPathID,
			testPathName: TestPathName,

			expectedResult: []byte(SUCCESSFUL_DELETE),
			expectedError:  false,
		},
		{
			name:         "Unsuccessful DeleteItemByID and no pathName",
			mockClient:   mockClientDeleteByIDErr(),
			testID:       TestID,
			testPathID:   string(TestInvalidPathID),
			testPathName: "",

			expectedResult: nil,
			expectedError:  true,
		},
		{
			name:         "Unsuccessful DeleteItemByID has pathname",
			mockClient:   mockClientDeleteByIDErr(),
			testID:       TestID,
			testPathID:   TestPathID,
			testPathName: TestPathName,

			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client = test.mockClient
			actual, err := DeleteItem(test.testID, test.testPathID, test.testPathName, TestPort)

			if test.expectedError && err == nil {
				t.Error("Expected an error")
				return
			}

			if !test.expectedError && err != nil {
				t.Errorf("Unexpectedly encountered error: %s", err)
				return
			}

			if !reflect.DeepEqual(test.expectedResult, actual) {
				t.Errorf("Expected result does not match the observed.\nExpected: %v\nObserved: %v\n", test.expectedResult, actual)
				return
			}
		})
	}
}
