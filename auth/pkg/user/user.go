package user

type User struct {
	Id 		 string

}

type RequestBucket struct {
	 maxrequest int64
	 Id    	      string
   Ip           string
	 Count        int64
}

func Verifyuser(id string) bool {
	//it will verify user using id. id has an hash key if exists return true else false
	
}

func Verifyhash(hash string) bool {
	//it will verify user using hash key. if exists return true else false
	//hash key is unique for each user.and it has an secret key
	
}

type Bucketbuffer struct {
	store sync.Map[string]RequestBucket,
	expiry sync.Map[string]int32
	
}
