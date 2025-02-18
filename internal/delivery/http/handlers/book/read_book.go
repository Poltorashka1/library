package bookhandlers

// todo add validation query params

//func (h *bookHandlers) ReadBook(w http.ResponseWriter, r *http.Request) {
//	// todo add other file format
//	//parts := strings.Split(r.URL.Path, "/")
//	//
//	//bookUUID := parts[2]
//	//var bookChapter string
//	//if len(parts) == 3 {
//	//	bookChapter = "1"
//	//} else {
//	//	bookChapter = parts[2]
//	//}
//
//	bookUUID := request.URLParse(r, "uuid")
//	bookChapter := request.URLParse(r, "chapter")
//	if bookChapter == "" {
//		bookChapter = "1"
//	}
//
//	var payload = dtos.BookFileRequest{
//		FileName: bookUUID,
//		FileType: h.cfg.FormatHTML(),
//		Chapter:  bookChapter,
//	}
//
//	result, err := h.useCase.BookFile(r.Context(), payload)
//	if err != nil {
//		switch {
//		case errors.Is(err, apperrors.ErrBookNotExist):
//			response.Error(w, err, http.StatusNotFound)
//			//response.Redirect(w, r, h.cfg.NotFoundURL())
//		default:
//			response.Error(w, err, http.StatusInternalServerError)
//		}
//		return
//	}
//
//	response.Success(w, result.File, h.cfg.HTML())
//}

//func getChapter(fileName string) (string, []byte, error) {
//	z, err := zip.OpenReader(fileName)
//	if err != nil {
//		return "", nil, err
//	}
//	defer z.Close()
//
//	for _, f := range z.File {
//		if strings.HasSuffix(f.Name, ".html") {
//			rc, err := f.Open()
//			if err != nil {
//				return "", nil, err
//			}
//			defer rc.Close()
//
//			var content bytes.Buffer
//			if _, err := io.Copy(&content, rc); err != nil {
//				return "", nil, err
//			}
//
//			return f.Name, content.Bytes(), nil
//		}
//	}
//
//	return "", nil, nil
//}
