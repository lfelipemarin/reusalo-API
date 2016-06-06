package main

import (
"github.com/zabawaba99/firego"
"github.com/gin-gonic/gin"
"github.com/itsjamie/gin-cors"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"os"
//"strconv"
"fmt"
"time"
)

// Direccion de la app en Firebase (cambiar en produccion por la del profe)
const APP_URL = "https://fir-android-b69c3.firebaseio.com"
const MONGO_URL = "mongodb://android:root@ds023902.mlab.com:23902/reusalo_db"

// Estructura Producto
type Producto struct {
	IdProd 		int `bson:"id_prod" json:"id_prod"`
	IdUsuario   int `bson:"id_usuario" json:"id_usuario"`
	NombreProd  string `bson:"nombre_prod" json:"nombre_prod"`
	FotoProd	string `bson:"foto_prod" json:"foto_prod"`
	DescrProd	string `bson:"descripcion_prod" json:"descripcion_prod"`
	Estado		string `bson:"estado" json:"estado"`
	FechaPub	string `bson:"fecha_publicacion" json:"fecha_publicacion"`
}

// Estructura Cliente
type Categorias struct {
	Id              bson.ObjectId `bson:"_id" json:"_id"`
	Categoria 		[]Categoria `bson:"categorias" json:"categorias"`
	//IdCategoria  	int `bson:"id_categoria" json:"id_categoria"`
	//NombreCat  		string `bson:"nombre" json:"nombre"`
	//Productos       []Producto `bson:"productos" json:"productos"`
}

type Categoria struct{
	Id              bson.ObjectId `bson:"id" json:"id"`
	IdCategoria  	int `bson:"id_categoria" json:"id_categoria"`
	NombreCat  		string `bson:"nombre" json:"nombre"`
	Productos       []Producto `bson:"productos" json:"productos"`
}

// Define las rutas de la API y la ejecuta
func main() {
	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge: 50 * time.Second,
		Credentials: true,
		ValidateHeaders: false,
		}))
	v1 := r.Group("")
	{
		v1.GET("/", ImOk)
		v1.GET("/categorias", GetCategorias)
		//v1.GET("/productos/:token", GetProductos)
		// v1.GET("/cliente/documento/:documento/:token", GetCliente)
		//v1.GET("/cliente/:correo/:token", GetClientePorCorreo)
		//v1.GET("/ejecutivo/:correo/:token", GetEjecutivoPorCorreoCliente)
		// v1.POST("/usuarios", PostUser)
		// v1.PUT("/usuarios/:id", UpdateUser)
		// v1.DELETE("/usuarios/:id", DeleteUser)
	}

	/*para correr en un puerto local*/
	//r.Run(":1337")
	r.Run("192.168.0.112:8081")
}

// Conecta a la base de datos
// Definir variable de entorno en Heroku: Settings, Config Vars
// Definir variable de entorno local: echo "export MONGO_URL=mongodb://goteam:goteam@ds019471.mlab.com:19471/my_bank_db" >> ~/.bashrc
func connect() (session *mgo.Session) {
	//connectURL := os.Getenv("MONGO_URL")
	connectURL := MONGO_URL
	fmt.Println(connectURL)
	session, err := mgo.Dial(connectURL)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	session.SetSafe(&mgo.Safe{})
	return session
}

// Autentica un token de Google en Firebase
// Argumentos:
// token: token generado por Google
func auth(token string) bool {
	f := firego.New(APP_URL, nil)
	f.Auth(token)
	return processResponse(f)
}

// Procesa la respuesta de una peticion a Firebase
// Argumentos:
// f: instancia de Firebase con que se hizo peticion
func processResponse(f *firego.Firebase) bool {
	var v map[string]interface{}
	if err := f.Value(&v); err != nil {
		return false
	}
	return true
}

func ImOk(ginContext *gin.Context) {
	ginContext.JSON(200, gin.H{
		"status":  "Im OK!!!",
		})
}

func GetCategorias(ginContext *gin.Context) {
	//token := ginContext.Params.ByName("token")
	//if auth(token){
		session := connect();
		defer session.Close()
		collection := session.DB("reusalo_db").C("categorias")

		categorias := []Categorias{}
		err := collection.Find(nil).All(&categorias)
		if err != nil {
			panic(err)
		}
		ginContext.JSON(200, categorias)
	//}else{
	//	ginContext.JSON(404, gin.H{
	//		"error":  "permiso denegado",
	//		})
	//}
}


