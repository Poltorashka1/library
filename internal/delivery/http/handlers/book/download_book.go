package bookhandlers

//func (h *bookHandlers) DownloadBook(w http.ResponseWriter, r *http.Request) {
//	// todo add other file format
//	parts := strings.Split(r.URL.Path, "/")
//	bookUUID := parts[len(parts)-1]
//
//	query := r.URL.Query()
//	bookFormat := query.Get("format")
//	if bookFormat == "" {
//		bookFormat = h.cfg.FormatHTML()
//	}
//
//	bookChapter := query.Get("chapter")
//	if bookChapter == "" && bookFormat == h.cfg.FormatHTML() {
//		bookChapter = "1"
//	}
//
//	var payload = dtos.BookFileRequest{
//		FileName: bookUUID,
//		FileType: bookFormat,
//		Chapter:  bookChapter,
//	}
//
//	result, err := h.useCase.BookFile(r.Context(), payload)
//	if err != nil {
//		// todo return
//		//switch {
//		//case errors.Is(err, apperrors.ErrBookNotExist):
//		//	response.Redirect(w, r, h.config.NotFoundURL())
//		//}
//		response.Error(w, err, http.StatusInternalServerError)
//		return
//	}
//
//	switch payload.FileType {
//	case h.cfg.FormatPDF():
//		response.Success(w, result.File, h.cfg.PDF(), map[string]string{"Content-Disposition": "attachment"}) // inline to read in browser
//	case h.cfg.FormatJSON():
//		response.Success(w, result.File, h.cfg.JSON(), nil)
//	default:
//		response.Error(w, err, http.StatusInternalServerError)
//		return
//	}
//
//	return
//	//w.Header().Set("Content-Type", h.config.PDF())
//	//
//	//if _, err := io.Copy(w, result.File); err != nil {
//	//	http.Error(w, "failed to send file", http.StatusInternalServerError)
//	//	return
//	//}
//	//response.Success(w, result.File, "application/html")
//	//_, b, err := getChapter(book)
//	//if err != nil {
//	//	http.Error(w, err.Error(), http.StatusInternalServerError)
//	//	return
//	//}
//	//w.Header().Set("Content-Type", "application/html")
//	//w.WriteHeader(http.StatusOK)
//	//w.Write(b)
//	////http.ServeFile(w, r, book)
//
//}
