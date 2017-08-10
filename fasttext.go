package fasttextgo

// #cgo LDFLAGS: -L${SRCDIR} -lfasttext -lstdc++
// #include <stdlib.h>
// void load_model(char *path);
// int predict(char *query, float *prob, char *buf, int buf_sz);
// char **predictK(char *query, int k, float **probs, int buf_sz);
import "C"
import (
	"errors"
	"unsafe"
)

// LoadModel - load FastText model
func LoadModel(path string) {
	C.load_model(C.CString(path))
}

// Predict - predict
func Predict(sentence string) (prob float32, label string, err error) {
	var cprob C.float
	var buf *C.char
	buf = (*C.char)(C.calloc(64, 1))

	ret := C.predict(C.CString(sentence), &cprob, buf, 64)

	if ret != 0 {
		err = errors.New("error in prediction")
	} else {
		label = C.GoString(buf)
		prob = float32(cprob)
	}
	C.free(unsafe.Pointer(buf))

	return prob, label, err
}

// PredictK returns top K predictions
func PredictK(sentence string, k int) (probs []float32, labels []string, err error) {
	// Make probls and labels to return
	probs = make([]float32, k)
	labels = make([]string, k)

	p := C.predictK(C.CString(sentence), C.int(k), (**C.float)(unsafe.Pointer(&probs[0])), 64)
	// TODO: How do we check for null pointer being returned

	q := ((*[1 << 30]*C.char)(unsafe.Pointer(p)))[:k]
	for i, cs := range q {
		labels[i] = C.GoString(cs)
		C.free(unsafe.Pointer(cs))
	}
	C.free(unsafe.Pointer(p))

	return probs, labels, nil
}