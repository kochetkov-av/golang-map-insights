package main

import (
	"fmt"
	"unsafe"
)

func showSpread(m interface{}) {
	// dataOffset is where the cell data begins in a bmap
	const dataOffset = unsafe.Offsetof(struct {
		tophash [bucketCnt]uint8
		cells   int64
	}{}.cells)

	t, h := mapTypeAndValue(m)

	fmt.Printf("Overflow buckets: %d", h.noverflow)

	numBuckets := 1 << h.B

	for r := 0; r < numBuckets*bucketCnt; r++ {
		bucketIndex := r / bucketCnt
		cellIndex := r % bucketCnt

		if cellIndex == 0 {
			fmt.Printf("\nBucket %3d:", bucketIndex)
		}

		// lookup cell
		b := (*bmap)(add(h.buckets, uintptr(bucketIndex)*uintptr(t.bucketsize)))
		if b.tophash[cellIndex] == 0 {
			// cell is empty
			continue
		}

		k := add(unsafe.Pointer(b), dataOffset+uintptr(cellIndex)*uintptr(t.keysize))

		ei := emptyInterface{
			_type: unsafe.Pointer(t.key),
			value: k,
		}
		key := *(*interface{})(unsafe.Pointer(&ei))
		fmt.Printf(" %3d", key.(int))
	}
	fmt.Printf("\n\n")
}

func main() {
	m := make(map[int]int)

	for i := 0; i < 50; i++ {
		m[i] = i * 3
	}

	showSpread(m)

	m = make(map[int]int)

	for i := 0; i < 8; i++ {
		m[i] = i * 3
	}

	showSpread(m)
}

type emptyInterface struct {
	_type unsafe.Pointer
	value unsafe.Pointer
}

func mapTypeAndValue(m interface{}) (*maptype, *hmap) {
	ei := (*emptyInterface)(unsafe.Pointer(&m))
	return (*maptype)(ei._type), (*hmap)(ei.value)
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

const (
	// Maximum number of key/elem pairs a bucket can hold.
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits
)

// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}

// mapextra holds fields that are not present on all maps.
type mapextra struct {
	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap
}

// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

type maptype struct {
	typ        _type
	key        *_type
	elem       *_type
	bucket     *_type // internal type representing a hash bucket
	keysize    uint8  // size of key slot
	elemsize   uint8  // size of elem slot
	bucketsize uint16 // size of bucket
	flags      uint32
}

// typeAlg is also copied/used in reflect/type.go.
// keep them in sync.
type typeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type tflag uint8
type nameOff int32
type typeOff int32

type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
