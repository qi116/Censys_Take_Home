package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const url = "http://localhost:8080"

func main() {
	fmt.Println("Starting tests")

	test_getEmptyKey()
	test_setEmptyKey()
	test_setEmptyValue()
	test_deleteDoesNotExist()
	test_getDoesNotExist()
	test_getAndSetAndGetAndSetAndGetAndDeleteAndDelete()
}

func handleResponse(resp *http.Response, err error) map[string]interface{} {
	if err != nil {
		// fmt.Println("Error in request:", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		// fmt.Println("Non-OK HTTP status:", resp.Status)
		return nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Response body:", string(body))

	marshalErr := json.Unmarshal(body, &result)
	if marshalErr != nil {
		fmt.Println("Error parsing response JSON:", marshalErr)
		return nil
	}
	return result
}

/*
Helper function to test getValue and compare response value with expectedValue
*/
func testGet(key string, expectedValue string) bool {
	resp, err := http.Get(url + "/getValue/" + key)

	result := handleResponse(resp, err)
	if result == nil {
		return false // error printed in handleResponse()
	}

	if result["value"] != expectedValue {
		fmt.Printf("Get test failed for key %s: expected value %s, got %s\n", key, expectedValue, result["value"])
		return false
	}

	return true
}

/*
Helper function to test setValue and compare response value with expectedValue
*/
func testSet(key string, value string, expectedValue string) bool {
	body, _ := json.Marshal(map[string]string{
		"key":   key,
		"value": value,
	})
	resp, err := http.Post(url+"/setValue", "application/json", bytes.NewBuffer(body))

	result := handleResponse(resp, err)
	if result == nil {
		return false // error printed in handleResponse()
	}

	if result["response"] != expectedValue {
		fmt.Printf("Set test failed for key %s: expected value %s, got %s\n", key, expectedValue, result["response"])
		return false
	}
	return true

}

/*
Helper function to test deleteValue and compare response value with expectedValue
*/
func testDelete(key string, expectedValue string) bool {
	req, _ := http.NewRequest("DELETE", url+"/deleteValue/"+key, nil)
	resp, err := http.DefaultClient.Do(req)

	result := handleResponse(resp, err)
	if result == nil {
		return false // error printed in handleResponse()
	}

	if result["response"] != expectedValue {
		fmt.Printf("Delete test failed for key %s: expected value %s, got %s\n", key, expectedValue, result["value"])
		return false
	}
	return true
}

func test_getEmptyKey() {
	result := testGet("", "")
	if !result { // test should fail
		fmt.Println("test_getEmptyKey PASSED")
	} else {
		fmt.Println("test_getEmptyKey FAILED")
	}
}

func test_setEmptyKey() {
	result := testSet("", "someValue", "")
	if !result { // test should fail
		fmt.Println("test_getEmptyKey PASSED")
	} else {
		fmt.Println("test_getEmptyKey FAILED")
	}
}

func test_setEmptyValue() {
	result := testSet("testKey", "", "")
	if !result { // test should fail
		fmt.Println("test_getEmptyValue PASSED")
	} else {
		fmt.Println("test_getEmptyValue FAILED")
	}
}

func test_deleteDoesNotExist() {
	result := testDelete("nonExistentKey", "Key nonExistentKey does not exist")
	if result { // test should fail
		fmt.Println("test_deleteDoesNotExist PASSED")
	} else {
		fmt.Println("test_deleteDoesNotExist FAILED")
	}
}

func test_getDoesNotExist() {
	result := testGet("nonExistentKey", "")
	if result { // test should fail
		fmt.Println("test_getDoesNotExist PASSED")
	} else {
		fmt.Println("test_getDoesNotExist FAILED")
	}
}

// Note that doing the && leads to short circuiting, so tests will stop at first failure.
func test_getAndSetAndGetAndSetAndGetAndDeleteAndDelete() {
	result := testGet("testKey1", "")
	result = result && testSet("testKey1", "testValue1", "New key testKey1 added with value testValue1")
	result = result && testGet("testKey1", "testValue1")
	result = result && testSet("testKey1", "testValue2", "Key testKey1 existed with value testValue1. Now updated with value testValue2")
	result = result && testGet("testKey1", "testValue2")
	result = result && testDelete("testKey1", "Key testKey1 existed and is now deleted")
	result = result && testDelete("testKey1", "Key testKey1 does not exist")
	if result {
		fmt.Println("test_getAndSetAndGetAndSetAndGetAndDeleteAndDelete PASSED")
	} else {
		fmt.Println("test_getAndSetAndGetAndSetAndGetAndDeleteAndDelete FAILED")
	}
}
