// Package unityweb is a library for parsing, unpacking, and repacking Unity Web
// data files.
//
// # Unpacking a Unity Web data file into a directory
//
//	pkg, err := unityweb.FromPackageFile("/path/to/test.data")
//
//	if err != nil {
//		panic(err)
//	}
//
//	err = pkg.Dump("/path/to/output/directory")
//
//	if err != nil {
//		panic(err)
//	}
//
// # Packing a directory into a Unity Web data file
//
//	pkg, err := unityweb.PackDirectory("/path/to/input/directory")
//
//	if err != nil {
//		panic(err)
//	}
//
//	err = pkg.PackToFile("/path/to/output.data")
//
//	if err != nil {
//		panic(err)
//	}
package unityweb
