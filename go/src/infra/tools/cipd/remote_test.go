// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cipd

import (
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	. "infra/tools/cipd/common"
)

func TestRemoteImpl(t *testing.T) {
	mockInitiateUpload := func(c C, reply string) (*UploadSession, error) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/upload/SHA1/abc",
				Reply:  reply,
			},
		})
		return remote.initiateUpload("abc")
	}

	mockFinalizeUpload := func(c C, reply string) (bool, error) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/finalize/abc",
				Reply:  reply,
			},
		})
		return remote.finalizeUpload("abc")
	}

	mockRegisterInstance := func(c C, reply string) (*registerInstanceResponse, error) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/repo/v1/instance",
				Query: url.Values{
					"package_name": []string{"pkgname"},
					"instance_id":  []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
				},
				Reply: reply,
			},
		})
		return remote.registerInstance(Pin{"pkgname", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	}

	mockFetchInstance := func(c C, reply string) (*fetchInstanceResponse, error) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "GET",
				Path:   "/_ah/api/repo/v1/instance",
				Query: url.Values{
					"package_name": []string{"pkgname"},
					"instance_id":  []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
				},
				Reply: reply,
			},
		})
		return remote.fetchInstance(Pin{"pkgname", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	}

	mockFetchACL := func(c C, reply string) ([]PackageACL, error) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "GET",
				Path:   "/_ah/api/repo/v1/acl",
				Query:  url.Values{"package_path": []string{"pkgname"}},
				Reply:  reply,
			},
		})
		return remote.fetchACL("pkgname")
	}

	mockModifyACL := func(c C, changes []PackageACLChange, body, reply string) error {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/repo/v1/acl",
				Query:  url.Values{"package_path": []string{"pkgname"}},
				Body:   body,
				Reply:  reply,
			},
		})
		return remote.modifyACL("pkgname", changes)
	}

	mockSetRef := func(c C, reply string) error {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/repo/v1/ref",
				Query: url.Values{
					"package_name": []string{"pkgname"},
					"ref":          []string{"some-ref"},
				},
				Body:  `{"instance_id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`,
				Reply: reply,
			},
		})
		return remote.setRef("some-ref", Pin{"pkgname", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	}

	mockAttachTags := func(c C, tags []string, body, reply string) error {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/repo/v1/tags",
				Query: url.Values{
					"package_name": []string{"pkgname"},
					"instance_id":  []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
				},
				Body:  body,
				Reply: reply,
			},
		})
		return remote.attachTags(Pin{"pkgname", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, tags)
	}

	mockResolveVersion := func(c C, reply string) (Pin, error) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "GET",
				Path:   "/_ah/api/repo/v1/instance/resolve",
				Query: url.Values{
					"package_name": []string{"pkgname"},
					"version":      []string{"tag_key:value"},
				},
				Reply: reply,
			},
		})
		return remote.resolveVersion("pkgname", "tag_key:value")
	}

	Convey("makeRequest POST works", t, func(c C) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/method",
				Reply:  `{"value":"123"}`,
			},
		})
		var reply struct {
			Value string `json:"value"`
		}
		err := remote.makeRequest("cas/v1/method", "POST", nil, &reply)
		So(err, ShouldBeNil)
		So(reply.Value, ShouldEqual, "123")
	})

	Convey("makeRequest GET works", t, func(c C) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "GET",
				Path:   "/_ah/api/cas/v1/method",
				Reply:  `{"value":"123"}`,
			},
		})
		var reply struct {
			Value string `json:"value"`
		}
		err := remote.makeRequest("cas/v1/method", "GET", nil, &reply)
		So(err, ShouldBeNil)
		So(reply.Value, ShouldEqual, "123")
	})

	Convey("makeRequest handles fatal error", t, func(c C) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/method",
				Status: 403,
			},
		})
		var reply struct{}
		err := remote.makeRequest("cas/v1/method", "POST", nil, &reply)
		So(err, ShouldNotBeNil)
	})

	Convey("makeRequest handles retries", t, func(c C) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/method",
				Status: 500,
			},
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/method",
				Reply:  `{}`,
			},
		})
		var reply struct{}
		err := remote.makeRequest("cas/v1/method", "POST", nil, &reply)
		So(err, ShouldBeNil)
	})

	Convey("makeRequest gives up trying", t, func(c C) {
		calls := []expectedHTTPCall{}
		for i := 0; i < remoteMaxRetries; i++ {
			calls = append(calls, expectedHTTPCall{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/method",
				Status: 500,
			})
		}
		remote := mockRemoteImpl(c, calls)
		var reply struct{}
		err := remote.makeRequest("cas/v1/method", "POST", nil, &reply)
		So(err, ShouldNotBeNil)
	})

	Convey("initiateUpload ALREADY_UPLOADED", t, func(c C) {
		s, err := mockInitiateUpload(c, `{"status":"ALREADY_UPLOADED"}`)
		So(err, ShouldBeNil)
		So(s, ShouldBeNil)
	})

	Convey("initiateUpload SUCCESS", t, func(c C) {
		s, err := mockInitiateUpload(c, `{"status":"SUCCESS","upload_session_id":"123","upload_url":"http://localhost"}`)
		So(err, ShouldBeNil)
		So(s, ShouldResemble, &UploadSession{"123", "http://localhost"})
	})

	Convey("initiateUpload ERROR", t, func(c C) {
		s, err := mockInitiateUpload(c, `{"status":"ERROR","error_message":"boo"}`)
		So(err, ShouldNotBeNil)
		So(s, ShouldBeNil)
	})

	Convey("initiateUpload unknown status", t, func(c C) {
		s, err := mockInitiateUpload(c, `{"status":"???"}`)
		So(err, ShouldNotBeNil)
		So(s, ShouldBeNil)
	})

	Convey("initiateUpload bad reply", t, func(c C) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/upload/SHA1/abc",
				Status: 403,
			},
		})
		s, err := remote.initiateUpload("abc")
		So(err, ShouldNotBeNil)
		So(s, ShouldBeNil)
	})

	Convey("finalizeUpload MISSING", t, func(c C) {
		finished, err := mockFinalizeUpload(c, `{"status":"MISSING"}`)
		So(err, ShouldNotBeNil)
		So(finished, ShouldBeFalse)
	})

	Convey("finalizeUpload UPLOADING", t, func(c C) {
		finished, err := mockFinalizeUpload(c, `{"status":"UPLOADING"}`)
		So(err, ShouldBeNil)
		So(finished, ShouldBeFalse)
	})

	Convey("finalizeUpload VERIFYING", t, func(c C) {
		finished, err := mockFinalizeUpload(c, `{"status":"VERIFYING"}`)
		So(err, ShouldBeNil)
		So(finished, ShouldBeFalse)
	})

	Convey("finalizeUpload PUBLISHED", t, func(c C) {
		finished, err := mockFinalizeUpload(c, `{"status":"PUBLISHED"}`)
		So(err, ShouldBeNil)
		So(finished, ShouldBeTrue)
	})

	Convey("finalizeUpload ERROR", t, func(c C) {
		finished, err := mockFinalizeUpload(c, `{"status":"ERROR","error_message":"boo"}`)
		So(err, ShouldNotBeNil)
		So(finished, ShouldBeFalse)
	})

	Convey("finalizeUpload unknown status", t, func(c C) {
		finished, err := mockFinalizeUpload(c, `{"status":"???"}`)
		So(err, ShouldNotBeNil)
		So(finished, ShouldBeFalse)
	})

	Convey("finalizeUpload bad reply", t, func(c C) {
		remote := mockRemoteImpl(c, []expectedHTTPCall{
			{
				Method: "POST",
				Path:   "/_ah/api/cas/v1/finalize/abc",
				Status: 403,
			},
		})
		finished, err := remote.finalizeUpload("abc")
		So(err, ShouldNotBeNil)
		So(finished, ShouldBeFalse)
	})

	Convey("registerInstance REGISTERED", t, func(c C) {
		result, err := mockRegisterInstance(c, `{
				"status": "REGISTERED",
				"instance": {
					"registered_by": "user:abc@example.com",
					"registered_ts": "1420244414571500"
				}
			}`)
		So(err, ShouldBeNil)
		So(result, ShouldResemble, &registerInstanceResponse{
			registeredBy: "user:abc@example.com",
			registeredTs: time.Unix(0, 1420244414571500000),
		})
	})

	Convey("registerInstance ALREADY_REGISTERED", t, func(c C) {
		result, err := mockRegisterInstance(c, `{
				"status": "ALREADY_REGISTERED",
				"instance": {
					"registered_by": "user:abc@example.com",
					"registered_ts": "1420244414571500"
				}
			}`)
		So(err, ShouldBeNil)
		So(result, ShouldResemble, &registerInstanceResponse{
			alreadyRegistered: true,
			registeredBy:      "user:abc@example.com",
			registeredTs:      time.Unix(0, 1420244414571500000),
		})
	})

	Convey("registerInstance UPLOAD_FIRST", t, func(c C) {
		result, err := mockRegisterInstance(c, `{
				"status": "UPLOAD_FIRST",
				"upload_session_id": "upload_session_id",
				"upload_url": "http://upload_url"
			}`)
		So(err, ShouldBeNil)
		So(result, ShouldResemble, &registerInstanceResponse{
			uploadSession: &UploadSession{"upload_session_id", "http://upload_url"},
		})
	})

	Convey("registerInstance ERROR", t, func(c C) {
		result, err := mockRegisterInstance(c, `{
				"status": "ERROR",
				"error_message": "Some error message"
			}`)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("registerInstance unknown status", t, func(c C) {
		result, err := mockRegisterInstance(c, `{"status":"???"}`)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("fetchInstance SUCCESS", t, func(c C) {
		result, err := mockFetchInstance(c, `{
				"status": "SUCCESS",
				"instance": {
					"registered_by": "user:abc@example.com",
					"registered_ts": "1420244414571500"
				},
				"fetch_url": "https://fetch_url"
			}`)
		So(err, ShouldBeNil)
		So(result, ShouldResemble, &fetchInstanceResponse{
			registeredBy: "user:abc@example.com",
			registeredTs: time.Unix(0, 1420244414571500000),
			fetchURL:     "https://fetch_url",
		})
	})

	Convey("fetchInstance PACKAGE_NOT_FOUND", t, func(c C) {
		result, err := mockFetchInstance(c, `{"status": "PACKAGE_NOT_FOUND"}`)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("fetchInstance INSTANCE_NOT_FOUND", t, func(c C) {
		result, err := mockFetchInstance(c, `{"status": "INSTANCE_NOT_FOUND"}`)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("fetchInstance ERROR", t, func(c C) {
		result, err := mockFetchInstance(c, `{
				"status": "ERROR",
				"error_message": "Some error message"
			}`)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("fetchACL SUCCESS", t, func(c C) {
		result, err := mockFetchACL(c, `{
				"status": "SUCCESS",
				"acls": {
					"acls": [
						{
							"package_path": "a",
							"role": "OWNER",
							"principals": ["user:a", "group:b"],
							"modified_by": "user:abc@example.com",
							"modified_ts": "1420244414571500"
						},
						{
							"package_path": "a/b",
							"role": "READER",
							"principals": ["group:c"],
							"modified_by": "user:abc@example.com",
							"modified_ts": "1420244414571500"
						}
					]
				}
			}`)
		So(err, ShouldBeNil)
		So(result, ShouldResemble, []PackageACL{
			{
				PackagePath: "a",
				Role:        "OWNER",
				Principals:  []string{"user:a", "group:b"},
				ModifiedBy:  "user:abc@example.com",
				ModifiedTs:  time.Unix(0, 1420244414571500000),
			},
			{
				PackagePath: "a/b",
				Role:        "READER",
				Principals:  []string{"group:c"},
				ModifiedBy:  "user:abc@example.com",
				ModifiedTs:  time.Unix(0, 1420244414571500000),
			},
		})
	})

	Convey("fetchACL ERROR", t, func(c C) {
		result, err := mockFetchACL(c, `{
				"status": "ERROR",
				"error_message": "Some error message"
			}`)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("modifyACL SUCCESS", t, func(c C) {
		expected := `{
				"changes": [
					{
						"action": "GRANT",
						"role": "OWNER",
						"principal": "user:a@example.com"
					},
					{
						"action": "REVOKE",
						"role": "READER",
						"principal": "user:b@example.com"
					}
				]
			}`
		// Strip " ", "\t" and "\n".
		expected = strings.Replace(expected, " ", "", -1)
		expected = strings.Replace(expected, "\n", "", -1)
		expected = strings.Replace(expected, "\t", "", -1)

		err := mockModifyACL(c, []PackageACLChange{
			{
				Action:    GrantRole,
				Role:      "OWNER",
				Principal: "user:a@example.com",
			},
			{
				Action:    RevokeRole,
				Role:      "READER",
				Principal: "user:b@example.com",
			},
		}, expected, `{"status":"SUCCESS"}`)
		So(err, ShouldBeNil)
	})

	Convey("modifyACL ERROR", t, func(c C) {
		err := mockModifyACL(c, []PackageACLChange{}, `{"changes":null}`, `{
				"status": "ERROR",
				"error_message": "Error message"
			}`)
		So(err, ShouldNotBeNil)
	})

	Convey("setRef SUCCESS", t, func(c C) {
		So(mockSetRef(c, `{"status":"SUCCESS"}`), ShouldBeNil)
	})

	Convey("setRef bad ref", t, func(c C) {
		err := mockRemoteImpl(c, nil).setRef(
			"BAD REF",
			Pin{"pkgname", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
		So(err, ShouldNotBeNil)
	})

	Convey("setRef PROCESSING_NOT_FINISHED_YET", t, func(c C) {
		err := mockSetRef(c, `{"status":"PROCESSING_NOT_FINISHED_YET", "error_message":"Blah"}`)
		So(err, ShouldResemble, &pendingProcessingError{message: "Blah"})
	})

	Convey("setRef ERROR", t, func(c C) {
		So(mockSetRef(c, `{"status":"ERROR", "error_message":"Blah"}`), ShouldNotBeNil)
	})

	Convey("attachTags SUCCESS", t, func(c C) {
		err := mockAttachTags(
			c, []string{"tag1:value1", "tag2:value2"},
			`{"tags":["tag1:value1","tag2:value2"]}`,
			`{"status":"SUCCESS"}`)
		So(err, ShouldBeNil)
	})

	Convey("attachTags bad tag", t, func(c C) {
		err := mockRemoteImpl(c, nil).attachTags(
			Pin{"pkgname", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			[]string{"BADTAG"})
		So(err, ShouldNotBeNil)
	})

	Convey("attachTags PROCESSING_NOT_FINISHED_YET", t, func(c C) {
		err := mockAttachTags(
			c, []string{"tag1:value1", "tag2:value2"},
			`{"tags":["tag1:value1","tag2:value2"]}`,
			`{"status":"PROCESSING_NOT_FINISHED_YET", "error_message":"Blah"}`)
		So(err, ShouldResemble, &pendingProcessingError{message: "Blah"})
	})

	Convey("attachTags ERROR", t, func(c C) {
		err := mockAttachTags(
			c, []string{"tag1:value1", "tag2:value2"},
			`{"tags":["tag1:value1","tag2:value2"]}`,
			`{"status":"ERROR", "error_message":"Blah"}`)
		So(err, ShouldNotBeNil)
	})

	Convey("resolveVersion SUCCESS", t, func(c C) {
		pin, err := mockResolveVersion(c, `{
			"status": "SUCCESS",
			"instance_id": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}`)
		So(err, ShouldBeNil)
		So(pin, ShouldResemble, Pin{"pkgname", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	})

	Convey("resolveVersion SUCCESS and bad instance ID", t, func(c C) {
		_, err := mockResolveVersion(c, `{
			"status": "SUCCESS",
			"instance_id": "bad_id"
		}`)
		So(err, ShouldNotBeNil)
	})

	Convey("resolveVersion PACKAGE_NOT_FOUND", t, func(c C) {
		_, err := mockResolveVersion(c, `{"status": "PACKAGE_NOT_FOUND"}`)
		So(err, ShouldNotBeNil)
	})

	Convey("resolveVersion INSTANCE_NOT_FOUND", t, func(c C) {
		_, err := mockResolveVersion(c, `{"status": "INSTANCE_NOT_FOUND"}`)
		So(err, ShouldNotBeNil)
	})

	Convey("resolveVersion AMBIGUOUS_VERSION", t, func(c C) {
		_, err := mockResolveVersion(c, `{"status": "AMBIGUOUS_VERSION"}`)
		So(err, ShouldNotBeNil)
	})

	Convey("resolveVersion ERROR", t, func(c C) {
		_, err := mockResolveVersion(c, `{"status": "ERROR", "error_message":"Blah"}`)
		So(err, ShouldNotBeNil)
	})

	Convey("resolveVersion bad status", t, func(c C) {
		_, err := mockResolveVersion(c, `{"status": "HUH?"}`)
		So(err, ShouldNotBeNil)
	})
}

////////////////////////////////////////////////////////////////////////////////

func mockRemoteImpl(c C, expectations []expectedHTTPCall) *remoteImpl {
	return &remoteImpl{mockClient(c, "", expectations)}
}
