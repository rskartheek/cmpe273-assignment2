package main

import (
	// Standard library packages
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// Third party packages
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Location represents the structure of our resource
	Location struct {
		Name       string        `json:"name" bson:"name"`
		Address    string        `json:"address" bson:"address"`
		City       string        `json:"city" bson:"city"`
		State      string        `json:"state" bson:"state"`
		Zip        string        `json:"zip" bson:"zip"`
		ID         bson.ObjectId `json:"id" bson:"_id"`
		Coordinate struct {
			Lat float32 `json:"latitude"`
			Lng float32 `json:"longitude"`
		} `json:"coordinate"`
	}

	//CoordinateResponse struct
	CoordinateResponse struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float32 `json:"lat"`
					Lng float32 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	//UpdateLocation represents the structure of location which should be updated
	UpdateLocation struct {
		Address    string `json:"address" bson:"address"`
		City       string `json:"city" bson:"city"`
		State      string `json:"state" bson:"state"`
		Zip        string `json:"zip" bson:"zip"`
		Coordinate struct {
			Lat float32 `json:"latitude"`
			Lng float32 `json:"longitude"`
		} `json:"coordinate"`
	}
)

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://assignment2:cmpe273Assignment2@ds037824.mongolab.com:37824/go_rest_api_assignment2")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

type (
	// LocationController represents the controller for operating on the User resource
	LocationController struct {
		session *mgo.Session
	}
)

//NewLocationController creates a new location controller
func NewLocationController(s *mgo.Session) *LocationController {
	return &LocationController{s}
}

// GetLocation retrieves an individual user resource
func (lc LocationController) GetLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub Location
	var locationObject Location

	// Fetch user
	if err := lc.session.DB("go_rest_api_assignment2").C("locations").FindId(oid).One(&locationObject); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(locationObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateLocation creates a new user resource
func (lc LocationController) CreateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub a location to be populated from the body

	var locationObject Location
	//String to store address
	var queryParamBuilder string
	// Populate the user data
	json.NewDecoder(r.Body).Decode(&locationObject)

	addressKeys := strings.Fields(locationObject.Address)
	cityKeys := strings.Fields(locationObject.City)
	stateKeys := strings.Fields(locationObject.State)
	keys := append(addressKeys, cityKeys...)
	locationKeys := append(keys, stateKeys...)
	for i := 0; i < len(locationKeys); i++ {
		if i == len(locationKeys)-1 {
			queryParamBuilder += locationKeys[i]
		} else {
			queryParamBuilder += locationKeys[i] + "+"
		}
	}
	url := fmt.Sprintf("http://maps.google.com/maps/api/geocode/json?address=%s&sensor=false", queryParamBuilder)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	// read json http response
	jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	var coordinates CoordinateResponse

	err = json.Unmarshal(jsonDataFromHTTP, &coordinates) // here!
	if err != nil {
		panic(err)
	}
	if len(coordinates.Results) == 0 {
		w.WriteHeader(400)
		return
	}
	locationObject.Coordinate.Lat = coordinates.Results[0].Geometry.Location.Lat
	locationObject.Coordinate.Lng = coordinates.Results[0].Geometry.Location.Lng

	// Add an Id
	locationObject.ID = bson.NewObjectId()

	// Write the user to mongo
	lc.session.DB("go_rest_api_assignment2").C("locations").Insert(locationObject)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(locationObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// ModifyLocation modifies a created resource
func (lc LocationController) ModifyLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)
	// Stub a location to be populated from the body

	var locationObject Location
	var retrivedObject Location

	//String to store address
	var queryParamBuilder string
	// Populate the user data
	json.NewDecoder(r.Body).Decode(&locationObject)

	addressKeys := strings.Fields(locationObject.Address)
	cityKeys := strings.Fields(locationObject.City)
	stateKeys := strings.Fields(locationObject.State)
	keys := append(addressKeys, cityKeys...)
	locationKeys := append(keys, stateKeys...)
	for i := 0; i < len(locationKeys); i++ {
		if i == len(locationKeys)-1 {
			queryParamBuilder += locationKeys[i]
		} else {
			queryParamBuilder += locationKeys[i] + "+"
		}
	}
	url := fmt.Sprintf("http://maps.google.com/maps/api/geocode/json?address=%s&sensor=false", queryParamBuilder)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	// read json http response
	jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	var coordinates CoordinateResponse

	err = json.Unmarshal(jsonDataFromHTTP, &coordinates) // here!
	if err != nil {
		panic(err)
	}
	if len(coordinates.Results) == 0 {
		w.WriteHeader(400)
		return
	}
	locationObject.Coordinate.Lat = coordinates.Results[0].Geometry.Location.Lat
	locationObject.Coordinate.Lng = coordinates.Results[0].Geometry.Location.Lng

	//Fetch user
	if err := lc.session.DB("go_rest_api_assignment2").C("locations").FindId(oid).One(&retrivedObject); err != nil {
		w.WriteHeader(404)
		return
	}
	locationObject.Name = retrivedObject.Name
	locationObject.ID = retrivedObject.ID
	if locationObject.City == "" {
		locationObject.City = retrivedObject.City
	}
	if locationObject.Address == "" {
		locationObject.Address = retrivedObject.Address
	}
	if locationObject.State == "" {
		locationObject.State = retrivedObject.State
	}
	if locationObject.Zip == "" {
		locationObject.Zip = retrivedObject.Zip
	}
	if err := lc.session.DB("go_rest_api_assignment2").C("locations").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}
	// Write the user to mongo
	lc.session.DB("go_rest_api_assignment2").C("locations").Insert(locationObject)

	// Marshal provided interface into JSON structure

	uj, _ := json.Marshal(locationObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)

}

// DeleteLocation removes an existing user resource
func (lc LocationController) DeleteLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := lc.session.DB("go_rest_api_assignment2").C("locations").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a LocationController instance
	lc := NewLocationController(getSession())

	// Get a location resource
	r.GET("/location/:id", lc.GetLocation)

	//Save a location
	r.POST("/location", lc.CreateLocation)

	//Update a Location
	r.PUT("/location/:id", lc.ModifyLocation)
	//Delete a location
	r.DELETE("/location/:id", lc.DeleteLocation)

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}
