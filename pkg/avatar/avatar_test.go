package avatar

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_RandomImage(t *testing.T) {
	Convey("Generate a random avatar from email", t, func() {
		_, err := RandomImage([]byte("gitote@local"))
		So(err, ShouldBeNil)

		Convey("Try to generate an image with size zero", func() {
			_, err := RandomImageSize(0, []byte("gitote@local"))
			So(err, ShouldNotBeNil)
		})
	})
}
