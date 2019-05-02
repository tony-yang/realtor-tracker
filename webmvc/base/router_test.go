package base_test

import (
	"testing"
	"webmvc/base"
	"webmvc/tester"
)

var routes *base.Routes

func setupSubTest(t *testing.T) func(t *testing.T) {
	t.Log("## Setup Sub Test")
	routes = base.CreateNewRouter()
	return func(t *testing.T) {
		t.Log("## Teardown Sub Test")
	}
}

func TestRegisterRoute(t *testing.T) {
	t.Run("Register resource only path", func(t *testing.T) {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		baseController := &base.Controller{}
		routes.RegisterRoute("/resources", baseController)

		tester.AssertTrue(t, routes.Exists("/resources"))
		expectedController, _ := routes.GetController("/resources")
		tester.AssertObjectEqual(t, expectedController, baseController)
	})

	t.Run("Register 2 resources only path", func(t *testing.T) {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		baseController := &base.Controller{}
		routes.RegisterRoute("/resources", baseController)
		routes.RegisterRoute("/resources2", baseController)

		tester.AssertTrue(t, routes.Exists("/resources"))
		tester.AssertTrue(t, routes.Exists("/resources2"))
		tester.AssertFalse(t, routes.Exists("/resources3"))
		expectedController1, _ := routes.GetController("/resources")
		expectedController2, _ := routes.GetController("/resources2")
		tester.AssertObjectEqual(t, expectedController1, baseController)
		tester.AssertObjectEqual(t, expectedController2, baseController)
	})

	t.Run("Register resources and pass resource id", func(t *testing.T) {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		baseController := &base.Controller{}
		routes.RegisterRoute("/resources", baseController)

		urlPath := "/resources/1"
		tester.AssertTrue(t, routes.Exists(urlPath))
		expectedController, expectedParam := routes.GetController(urlPath)
		tester.AssertObjectEqual(t, expectedController, baseController)
		tester.AssertStringEqual(t, expectedParam, "1")
	})

	t.Run("Register resources and pass in sub-resource", func(t *testing.T) {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		baseController := &base.Controller{}
		routes.RegisterRoute("/resources", baseController)

		urlPath := "/resources/1/resources2/2"
		tester.AssertTrue(t, routes.Exists(urlPath))
		expectedController, expectedParam := routes.GetController(urlPath)
		tester.AssertObjectEqual(t, expectedController, baseController)
		tester.AssertStringEqual(t, expectedParam, "1/resources2/2")
	})

	t.Run("Register resources and pass in parameters", func(t *testing.T) {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		baseController := &base.Controller{}
		routes.RegisterRoute("/resources", baseController)

		urlPath := "/resources/1?query=value"
		tester.AssertTrue(t, routes.Exists(urlPath))
		expectedController, expectedParam := routes.GetController(urlPath)
		tester.AssertObjectEqual(t, expectedController, baseController)
		tester.AssertStringEqual(t, expectedParam, "1?query=value")
	})
}
