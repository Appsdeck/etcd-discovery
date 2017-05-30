package service

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceAll(t *testing.T) {
	Convey("With no services", t, func() {
		s, err := Get("service-test-get-1")
		So(err, ShouldBeNil)

		hosts, err := s.All()
		So(err, ShouldBeNil)
		So(len(hosts), ShouldEqual, 0)
	})

	Convey("With two services", t, func() {
		log.Println("COUILLLE PUTAIN DE MERDE !!!!!!!!")
		host1 := genHost("test1")
		host2 := genHost("test2")
		_, c1 := Register("test-get-222", host1, make(chan struct{}))
		_, c2 := Register("test-get-222", host2, make(chan struct{}))

		<-c1
		<-c2

		s, err := Get("test-get-222")
		hosts, err := s.All()
		So(err, ShouldBeNil)
		So(len(hosts), ShouldEqual, 2)
		if hosts[0].PrivateHostname == "test1-private.dev" {
			So(hosts[1].PrivateHostname, ShouldEqual, "test2-private.dev")
		} else {
			So(hosts[1].PrivateHostname, ShouldEqual, "test1-private.dev")
			So(hosts[0].PrivateHostname, ShouldEqual, "test2-private.dev")
		}
	})
}

func TestServiceFirst(t *testing.T) {
	Convey("With no services", t, func() {
		s, err := Get("service-test-1")
		So(err, ShouldBeNil)
		host, err := s.First()
		So(err, ShouldNotBeNil)
		So(host, ShouldBeNil)
		So(err.Error(), ShouldEqual, "No host found for this service")
	})

	Convey("With a service", t, func() {
		host1 := genHost("test1")
		_, c := Register("test-truc", host1, make(chan struct{}))
		<-c

		s, err := Get("test-truc")
		So(err, ShouldBeNil)
		host, err := s.First()
		So(err, ShouldBeNil)
		So(host, ShouldNotBeNil)
		So(host.Name, ShouldEqual, host1.Name)
	})
}

func TestServiceOne(t *testing.T) {
	Convey("With no services", t, func() {
		s, err := Get("service-test-1")
		So(err, ShouldBeNil)
		host, err := s.One()
		So(err, ShouldNotBeNil)
		So(host, ShouldBeNil)
		So(err.Error(), ShouldEqual, "No host found for this service")
	})

	Convey("With a service", t, func() {
		host1 := genHost("test1")
		_, c := Register("test-truc", host1, make(chan struct{}))
		<-c

		s, err := Get("test-truc")
		So(err, ShouldBeNil)
		host, err := s.One()
		So(err, ShouldBeNil)
		So(host, ShouldNotBeNil)
		So(host.Name, ShouldEqual, host1.Name)
	})
}

func TestServiceUrl(t *testing.T) {
	Convey("With a public service", t, func() {
		Convey("With a service without any password", func() {
			host := genHost("test")
			host.User = ""
			host.Password = ""
			_, c := Register("service-url-1", host, make(chan struct{}))

			<-c

			s, err := Get("service-url-1")
			So(err, ShouldBeNil)
			url, err := s.Url("http", "/path")
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "http://public.dev:10000/path")
		})

		Convey("With a host with a password", func() {
			host := genHost("test")
			_, c := Register("service-url-3", host, make(chan struct{}))

			<-c

			s, err := Get("service-url-3")
			So(err, ShouldBeNil)
			url, err := s.Url("http", "/path")
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "http://user:password@public.dev:10000/path")
		})

		Convey("When the port does'nt exists", func() {
			host := genHost("test")
			_, c := Register("service-url-4", host, make(chan struct{}))

			<-c

			s, err := Get("service-url-4")
			So(err, ShouldBeNil)
			url, err := s.Url("htjp", "/path")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unknown scheme")
			So(len(url), ShouldEqual, 0)
		})
	})
}