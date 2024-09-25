/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
 **/

package bitmap

import "testing"

func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(5)

	b.Set("pppp")
	b.Set("22")
	b.Set("ccc")
	b.Set("eee")

	for _, bit := range b.bits {
		t.Logf("%b %v", bit, bit)
	}
}
