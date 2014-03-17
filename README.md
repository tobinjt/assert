assert
======

A helper library to make writing Golang tests easier.

Are you fed up writing tests like this?

    import "reflect"
    import "testing"
    func TestSomething(t *testing.T) {
            result := something()
            expected := []int{2, 9, 6}
            if !reflect.DeepEqual(expected, result) {
                    t.Errorf("something(): %#v != %#v\n", expected, result)
            }
    }

Would you prefer to write tests like this?

    import "github.com/tobinjt/assert"
    import "testing"
    func TestSomething(t *testing.T) {
            assert.Equal(t, "something()", []int{2, 9, 6}, something())
    }

This package makes it easy.

All functions return true if the test passes, and false if the test fails.  This
allows you to write tests like:

    func TestSomething(t *testing.T) {
            result, err := something()
            if assert.ErrIsNil(t, "something()", err) {
                    assert.Equal(t, "something()", 7, result)
            }
    }
