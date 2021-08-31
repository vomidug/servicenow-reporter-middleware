package main

import (
	"fmt"
	"testing"
)

type testarrs struct {
	name string
	oldIncs []Incident
	newIncs []Incident
	appeared []Incident
	disappeared []Incident
}

var tests = []testarrs{
	{ 
		name: "Nothing Happens",
		oldIncs: []Incident{ 
			{ 
				Number:"1543124",
			},
		},
		newIncs: []Incident{
			{
				Number:"1543124",
			},
		},
	},	
	{ 
		name:"One Disappeared",
		oldIncs: []Incident{ 
			{ 
				Number:"1543124",
			},
			{ 
				Number:"1543122",
			},
		},
		newIncs: []Incident{
			{
				Number:"1543124",
			},
		},
		appeared: []Incident{},
		disappeared: []Incident{
			{
				Number:"1543122",
			},
		},
	},
	{ 
		name: "One Appeared",
		oldIncs: []Incident{ 
			{ 
				Number:"1543124",
			},
		},
		newIncs: []Incident{
			{
				Number:"1543124",
			},
			{ 
				Number:"1543122",
			},
		},
		appeared: []Incident{
			{
				Number:"1543122",
			},
		},
		disappeared: []Incident{
		},
	},	
	{ 
		name:"One Appeared One Disappeared",
		oldIncs: []Incident{ 
			{ 
				Number:"1",
			},
			{
				Number:"2",
			},
			{
				Number:"3",
			},
		},
		newIncs: []Incident{
			{
				Number:"1",
			},
			{
				Number:"2",
			},
			{
				Number:"4",
			},
		},
		appeared: []Incident{
			{
				Number:"4",
			},
		},
		disappeared: []Incident{
			{
				Number:"3",
			},
		},
	},
	{ 
		name: "Three Appeared Three Disappeared",
		oldIncs: []Incident{ 
			{ 
				Number:"1",
			},
			{
				Number:"2",
			},
			{
				Number:"3",
			},
			{
				Number:"4",
			},
			{
				Number:"5",
			},
			{
				Number:"6",
			},
		},
		newIncs: []Incident{
			{
				Number:"1",
			},
			{
				Number:"2",
			},
			{
				Number:"4",
			},
			{
				Number:"7",
			},
			{
				Number:"8",
			},
		},
		appeared: []Incident{
			{
				Number:"4",
			},
			{
				Number:"7",
			},
			{
				Number:"8",
			},
		},
		disappeared: []Incident{
			{
				Number:"3",
			},
			{
				Number:"5",
			},
			{
				Number:"6",
			},
		},
	},
	{ 
		name:"Was nothing, three appeared",
		oldIncs: []Incident{ 
			{ 
			},
		},
		newIncs: []Incident{
			{
				Number:"3",
			},
			{
				Number:"5",
			},
			{
				Number:"6",
			},
		},
		appeared: []Incident{
			{
				Number:"3",
			},
			{
				Number:"5",
			},
			{
				Number:"6",
			},
		},
		disappeared: []Incident{

		},
	},
}

	/* template, useful for copypasting with y16y
	{ // single PAIR of arrays
		oldIncs: []Incident{ // single array of incidents
			{ // single Incident in Array
			},
		},
		newIncs: []Incident{
			{
			},
		},
		appeared: []Incident{

		},
		disappeared: []Incident{

		},
	},
	*/



func TestCompareincs (t *testing.T) {

	for _, pair := range tests {
//		fmt.Println("Testing: ")
//		fmt.Println("Old: ")
//		fmt.Println(pair.oldIncs)
//		fmt.Println("New: ")
//		fmt.Println(pair.newIncs)
		fmt.Println(pair.name)
		appeared, disappeared := compareincs(pair.oldIncs, pair.newIncs)
//		fmt.Println("")
//		fmt.Println("Expected: ")
//		fmt.Println("Appeared: ")
//		fmt.Println(pair.appeared)
//		fmt.Println("Disappeared: ")
//		fmt.Println(pair.disappeared)
//		fmt.Println("")
//		fmt.Println("Received")
//		fmt.Println("Appeared: ")
//		fmt.Println(appeared)
//		fmt.Println("Disappeared: ")
//		fmt.Println(disappeared)
//		fmt.Println("")
//		fmt.Println("Result: ")

		if (compareArrays(pair.disappeared, disappeared) && compareArrays(pair.appeared, appeared)) {
			fmt.Println("PASSED")
		} else {

			fmt.Println("NOT PASSED")
			t.Error(" NOT PASSED " + pair.name)
		}

		fmt.Println("")
	}
}

func compareArrays(valid []Incident, checkable []Incident) bool {

//	fmt.Println("Testing Array:")
//	fmt.Println(checkable)
//	fmt.Println("Model Array:")
//	fmt.Println(valid)


	for i := 0; i < len(checkable); i++ {

		if checkable[i] != (Incident{}) {

			var flag bool = false

//			fmt.Println("Checking if")
//			fmt.Println(checkable[i])
//			fmt.Println("Exists in")
//			fmt.Println(valid)

			for j := 0; j < len(valid); j++ {

//				fmt.Println("Comparing")
//				fmt.Println(checkable[i])
//				fmt.Println("With")
//				fmt.Println(valid[j])

				if compareSingle(checkable[i], valid[j]) {
					flag = true
				}
			}

			if !flag {
//				fmt.Println("\t\t\tError!!!")
//				fmt.Println(checkable[i])
//				fmt.Println("Does not exist in")
//				fmt.Println(valid)
				return false
			}
		}

	}
	
	return true

}
