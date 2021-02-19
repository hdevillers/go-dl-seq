package dataset

import (
	"errors"
	
	"gonum.org/v1/hdf5"
)

const (
	LabelMaxLength int = 50 // Max sequence ID (in Bytes)
)


// Create a complete dataset in a H5 file
func AddOhoDataset(f *hdf5.File, gn string, x []uint16, y []uint16, l []string, ns int, ls int) error {
	// Convert labels from string to bytes
	lb := make([]byte, ns*LabelMaxLength)
	for i := 0; i < ns; i++ {
		lsb := []byte(l[i])
		if len(lsb) >= LabelMaxLength {
			return errors.New("To long sequence ID.")
		}
		at := i * LabelMaxLength
		copy(lb[at:], lsb)
	}

	// Create the group
	grp, err := f.CreateGroup(gn)
	if err != nil {
		return err
	}

	// Prepare dims
	dx := []uint{uint(ns), 4, uint(ls)}
	dy := []uint{uint(ns)}
	dl := []uint{uint(ns), uint(LabelMaxLength)}

	// Prepare spaces
	sx, err := hdf5.CreateSimpleDataspace(dx, nil)
	if err != nil {
		return err
	}
	sy, err := hdf5.CreateSimpleDataspace(dy, nil)
	if err != nil {
		return err
	}
	sl, err := hdf5.CreateSimpleDataspace(dl, nil)
	if err != nil {
		return err
	}

	// Prepare data types
	tx, err := hdf5.NewDatatypeFromValue(x[0])
	if err != nil {
		return err
	}
	ty, err := hdf5.NewDatatypeFromValue(y[0])
	if err != nil {
		return err
	}
	tl, err := hdf5.NewDatatypeFromValue(lb[0])
	if err != nil {
		return err
	}

	// Create the different datasets and wrote them
	dsx, err := grp.CreateDataset("x", tx, sx)
	if err != nil {
		return err
	}
	err = dsx.Write(&x)
	if err != nil {
		return err
	}
	err = dsx.Close()
	if err != nil {
		return err
	}

	dsy, err := grp.CreateDataset("y", ty, sy)
	if err != nil {
		return err
	}
	err = dsy.Write(&y)
	if err != nil {
		return err
	}
	err = dsy.Close()
	if err != nil {
		return err
	}

	dsl, err := grp.CreateDataset("labels", tl, sl)
	if err != nil {
		return err
	}
	err = dsl.Write(&lb)
	if err != nil {
		return err
	}
	err = dsl.Close()
	if err != nil {
		return err
	}

	// Close the group
	err = grp.Close()

	return nil
}