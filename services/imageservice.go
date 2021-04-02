package services

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
	//
)

// ConferenceService will do dirty work
type ImageService struct {
}

func (srv *ImageService) Validation(contenttype string) bool {
	isvalid := false
	filetypes := [4]string{"image/gif", "image/png", "image/jpeg", "image/jpg"}
	for key, _ := range filetypes {
		if filetypes[key] == contenttype {
			fmt.Println("filetype", contenttype)
			isvalid = true
		}
	}

	return isvalid
}
func (repo *ImageService) WriteFileToDisk(folder, filename string, file []byte) error {
	// _, err := os.Stat(folder)
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		mkdirerr := os.MkdirAll(folder, 0755)
	// 		if mkdirerr != nil {
	// 			fmt.Println("folder not created")
	// 			panic(mkdirerr)
	// 		}
	// 	} else {
	// 		return err
	// 	}
	// }
	fmt.Println("folder path==", folder)
	localfile, err := os.Create(path.Join(folder, filename))
	if err != nil {
		return errors.New("file does not created")
		fmt.Println("file does not created", err)
	}
	fmt.Println("created file is ", localfile)

	resBytes, writeerr := localfile.Write(file)
	if writeerr != nil {
		fmt.Println("file write error ", writeerr)
		return writeerr
	}
	fmt.Println("respose of file write ", resBytes)
	defer localfile.Close()
	fileinfo, err := os.Stat(folder + "/" + filename)
	if err != nil {
		fmt.Println("The fileinfo err", err)
		return err

	}
	//get the size
	size := fileinfo.Size()
	sizeinkb := float64(size) / 1024
	sizeinMb := sizeinkb / 1024
	fmt.Println("size in mb", sizeinMb)
	if sizeinMb > float64(revel.Config.IntDefault("hq.image.size", 2)) {

		return errors.New("greater than 2mb")
	}
	imagefile, err := os.Open(folder + "/" + filename)
	if err != nil {
		fmt.Println("err imagefile", err)
		imagefile.Close()
		return err
	}
	fmt.Println("imagefile", imagefile)

	image, _, err := image.DecodeConfig(imagefile)
	if err != nil {
		fmt.Println("error in decoding for dimentions", err)
		imagefile.Close()
		return err
	}
	fmt.Println("DIMENTIONS", image.Height, image.Width, revel.Config.IntDefault("hq.image.width", 624), revel.Config.IntDefault("hq.image.height", 825))
	if image.Width > revel.Config.IntDefault("hq.image.width", 624) {
		fmt.Println("error in dimentions", "invalid dimentions")
		imagefile.Close()
		return errors.New("invalid dimentions")
	}
	if image.Height > revel.Config.IntDefault("hq.image.height", 825) {
		imagefile.Close()
		return errors.New("invalid dimentions")
	}
	return nil
}

func (srv *ImageService) CreateImage(ImageData viewmodels.ImageVMWrite) error {
	// confID, _ := uuid.FromString(ConfData.ID)

	ImageRepo := repositories.Images{}
	// conferencerepo := repositories.Conferences{}
	image := models.Image{
		Name:           ImageData.Name,
		BasicURL:       ImageData.BasicURL,
		ImageURLPrefix: ImageData.ImageURLPrefix,
		FolderPath:     ImageData.FolderPath,
		EntityId:       ImageData.EntityID,
		EntityType:     ImageData.EntityType,
		IsActive:       ImageData.IsActive,
		ImageCategory:  ImageData.ImageCategory,
	}
	// conference.ID = confID

	ImageError := ImageRepo.Create(&image)

	if ImageError != nil {
		fmt.Printf("create image returns error : %v\n", ImageError)
		return ImageError

	}
	return nil
}

func (srv *ImageService) Update(ImageData viewmodels.ImageVMWrite) error {
	// confID, _ := uuid.FromString(ConfData.ID)

	ImageRepo := repositories.Images{}
	// conferencerepo := repositories.Conferences{}
	image := models.Image{
		Name:           ImageData.Name,
		BasicURL:       ImageData.BasicURL,
		ImageURLPrefix: ImageData.ImageURLPrefix,
		FolderPath:     ImageData.FolderPath,
		EntityId:       ImageData.EntityID,
		EntityType:     ImageData.EntityType,
		IsActive:       ImageData.IsActive,
		ImageCategory:  ImageData.ImageCategory,
	}
	// conference.ID = confID

	ImageDB, ImageError := ImageRepo.Update(image)
	fmt.Println("here is img id", ImageDB.ID)
	if ImageError != nil {
		fmt.Printf("create image returns error : %v\n", ImageError)
		return ImageError

	}
	return nil
}
func (srv *ImageService) GetImage(Entityid string, entityType string, imageCategory string) (models.Image, error) {

	fmt.Println("here in service getbyentityid", Entityid)
	entityId, _ := uuid.FromString(Entityid)
	imgrepo := repositories.Images{}
	ImgDB, err := imgrepo.GetImage(entityId, entityType, imageCategory)
	if err != nil {
		fmt.Println("Error in get image service", Entityid)
		return ImgDB, err
	}
	return ImgDB, nil
}

func (repo *ImageService) WriteprofileToDisk(folder, filename string, file []byte) error {

	fmt.Println("folder path==", folder)
	localfile, err := os.Create(path.Join(folder, filename))
	if err != nil {
		return errors.New("file does not created")
		fmt.Println("file does not created", err)
	}
	fmt.Println("created file is ", localfile)

	resBytes, writeerr := localfile.Write(file)
	if writeerr != nil {
		fmt.Println("file write error ", writeerr)
		return writeerr
	}
	fmt.Println("respose of file write ", resBytes)
	defer localfile.Close()
	fileinfo, err := os.Stat(folder + "/" + filename)
	if err != nil {
		fmt.Println("The fileinfo err", err)
		return err

	}
	//get the size
	size := fileinfo.Size()
	sizeinkb := float64(size) / 1024
	sizeinMb := sizeinkb / 1024
	fmt.Println("size in mb", sizeinMb)
	if sizeinMb > float64(revel.Config.IntDefault("hq.ProfileImage.size", 1)) {

		return errors.New("greater than 2mb")
	}
	imagefile, err := os.Open(folder + "/" + filename)
	if err != nil {
		fmt.Println("err imagefile", err)
		imagefile.Close()
		return err
	}
	fmt.Println("imagefile", imagefile)

	image, _, err := image.DecodeConfig(imagefile)
	if err != nil {
		fmt.Println("error in decoding for dimentions", err)
		imagefile.Close()
		return err
	}
	fmt.Println("DIMENTIONS", image.Height, image.Width, revel.Config.IntDefault("hq.image.width", 624), revel.Config.IntDefault("hq.image.height", 825))
	if image.Width > revel.Config.IntDefault("hq.image.width", 624) {
		fmt.Println("error in dimentions", "invalid dimentions")
		imagefile.Close()
		return errors.New("invalid dimentions")
	}
	if image.Height > revel.Config.IntDefault("hq.image.height", 825) {
		imagefile.Close()
		return errors.New("invalid dimentions")
	}
	return nil
}
