// Forked from https://github.com/akhenakh/hunspellgo
package hunspell
// #cgo linux LDFLAGS: -lhunspell
// #cgo darwin LDFLAGS: -lhunspell-1.3 -L/usr/local/Cellar/hunspell/1.3.2/lib
// #cgo darwin CFLAGS: -I/usr/local/Cellar/hunspell/1.3.2/include
// #cgo freebsd CFLAGS: -I/usr/local/include
// #cgo freebsd LDFLAGS: -L/usr/local/lib -lhunspell-1.3
//
// #include <stdlib.h>
// #include <stdio.h>
// #include <hunspell/hunspell.h>
import "C"

import (

    "sync"
    "unsafe"
    "reflect"
    "runtime"
)

type Hunhandle struct {
    handle *C.Hunhandle
    lock   *sync.Mutex
}

func Hunspell(affpath string, dpath string) *Hunhandle {

    affpathcs := C.CString(affpath)
    defer C.free(unsafe.Pointer(affpathcs))

    dpathcs := C.CString(dpath)
    defer C.free(unsafe.Pointer(dpathcs))

    h := &Hunhandle{lock: new(sync.Mutex)}
    h.handle = C.Hunspell_create(affpathcs, dpathcs)

    runtime.SetFinalizer(h, func(handle *Hunhandle) {
        C.Hunspell_destroy(handle.handle)
        h.handle = nil
    })

    return h
}

func CArrayToString(c **C.char, l int) []string {
    
    s := []string{}

    hdr := reflect.SliceHeader{
            Data: uintptr(unsafe.Pointer(c)),
            Len:  l,
            Cap:  l,
    }

    for _, v := range *(*[]*C.char)(unsafe.Pointer(&hdr)) {
        s = append(s, C.GoString(v))
    }

    return s
}

func (handle *Hunhandle) Suggest(word string) []string {
    wordcs := C.CString(word)
    defer C.free(unsafe.Pointer(wordcs))

    var carray **C.char
    var length C.int
    handle.lock.Lock()
    length = C.Hunspell_suggest(handle.handle, &carray, wordcs)
    handle.lock.Unlock()

    words := CArrayToString(carray, int(length))

    C.Hunspell_free_list(handle.handle, &carray, length)
    return words
}

func (handle *Hunhandle) Add(word string) bool {

    cWord := C.CString(word)
    defer C.free(unsafe.Pointer(cWord))

    var r C.int
    r = C.Hunspell_add(handle.handle, cWord)

    if int(r) != 0 {
        return false
    }

    return true
}

func (handle *Hunhandle) Stem(word string) []string {
    wordcs := C.CString(word)
    defer C.free(unsafe.Pointer(wordcs))
    var carray **C.char
    var length C.int
    handle.lock.Lock()
    length = C.Hunspell_stem(handle.handle, &carray, wordcs)
    handle.lock.Unlock()

    words := CArrayToString(carray, int(length))

    C.Hunspell_free_list(handle.handle, &carray, length)
    return words
}

func (handle *Hunhandle) Spell(word string) bool {
    wordcs := C.CString(word)
    defer C.free(unsafe.Pointer(wordcs))
    handle.lock.Lock()
    res := C.Hunspell_spell(handle.handle, wordcs)
    handle.lock.Unlock()

    if int(res) == 0 {
        return false
    }
    return true
}

// Morphological generation by example
func (handle *Hunhandle) Generate(word string, example string) []string {
	wordcs := C.CString(word)
	examplecs := C.CString(example)
	defer C.free(unsafe.Pointer(wordcs))
	defer C.free(unsafe.Pointer(examplecs))
	var carray **C.char
	var length C.int
	handle.lock.Lock()
	length = C.Hunspell_generate(handle.handle, &carray, wordcs, examplecs)
	handle.lock.Unlock()

	words := CArrayToString(carray, int(length))

	C.Hunspell_free_list(handle.handle, &carray, length)
	return words
}
