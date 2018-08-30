package markup_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "gitote.com/gitote/gitote/pkg/markup"
)

func Test_Sanitizer(t *testing.T) {
	NewSanitizer()
	Convey("Sanitize HTML string and bytes", t, func() {
		testCases := []string{
			// Regular
			`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`, `<a href="http://www.google.com" rel="nofollow">Google</a>`,

			// Code highlighting class
			`<code class="random string"></code>`, `<code></code>`,
			`<code class="language-random ui tab active menu attached animating sidebar following bar center"></code>`, `<code></code>`,
			`<code class="language-go"></code>`, `<code class="language-go"></code>`,

			// Input checkbox
			`<input type="hidden">`, ``,
			`<input type="checkbox">`, `<input type="checkbox">`,
			`<input checked disabled autofocus>`, `<input checked="" disabled="">`,
		}

		for i := 0; i < len(testCases); i += 2 {
			So(Sanitize(testCases[i]), ShouldEqual, testCases[i+1])
			So(string(SanitizeBytes([]byte(testCases[i]))), ShouldEqual, testCases[i+1])
		}
	})
}
