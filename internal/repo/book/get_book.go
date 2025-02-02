package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	apperrors "book/internal/errors"
	"context"
	"strings"
)

// todo get image

// Book get book by uuid, err return apperrors.ErrBookNotFound
func (r *bookRepository) Book(ctx context.Context, uuid string) (*entities.Book, error) {
	query := storage.Query{
		QueryName: "get book by uuid",
		Query: `select t1.uuid, 
       			 t1.isbn, 
       			 t1.title, 
       			 t1.publication_year, 
       			 t1.description, 
       			 t1.books_file_uuid,
       			 t1.publisher,
       			 t3.uuid,
       			 t3.nickname,
       			 t3.name,
       			 t3.surname,
       			 t3.patronymic
			from books as t1
         		join book_author as t2 on t2.book_id = t1.id
         		join authors as t3 on t2.author_id = t3.id
			WHERE t1.uuid = $1;`,
		Args: []any{uuid},
	}

	rows, err := r.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	// todo optimize при нескольких авторах происходит перезаписывание информации о книге
	var book entities.Book
	if !rows.Next() {
		return nil, apperrors.ErrBookNotFound
	}

	for {
		var author entities.Author
		err := rows.Scan(
			&book.UUID,
			&book.ISBN,
			&book.Title,
			&book.PublicationYear,
			&book.Description,
			&book.BooksFileUUID,
			&book.Publisher,
			&author.UUID,
			&author.NickName,
			&author.Name,
			&author.Surname,
			&author.Patronymic,
		)
		if err != nil {
			return nil, err
		}
		book.BookAuthors.Authors = append(book.BookAuthors.Authors, author)

		if !rows.Next() {
			break
		}
	}

	return &book, err
}

// Books get books with authors, genres; err return apperrors.ErrPageNotFound
func (r *bookRepository) Books(ctx context.Context, filter entities.BookFilter) (*entities.Books, error) {
	// todo mb in cache

	query := storage.Query{QueryName: "get books count", Query: "select count(*) from books"}
	var booksCount int
	err := r.db.QueryRowContext(ctx, query).Scan(&booksCount)
	if err != nil {
		return nil, err
	}
	if filter.Start > booksCount {
		return nil, apperrors.ErrPageNotFound
	}

	query = storage.Query{
		QueryName: "select books with pagination",
		Query: `
				SELECT t1.id,t1.uuid,t1.isbn,t1.title,t1.publication_year,t1.description,
       				GROUP_CONCAT(t3.nickname, ', ') AS authors,
       				GROUP_CONCAT(t3.uuid, ', ')     AS authors_uuid
				FROM books AS t1
         			JOIN book_author AS t2 ON t2.book_id = t1.id
         			JOIN authors AS t3 ON t2.author_id = t3.id
				WHERE t1.id > $1 and t1.id <= $2
				GROUP BY t1.id
				ORDER BY t1.id`,
		Args: []any{filter.Start, filter.Stop},
	}

	rows, err := r.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var books entities.Books

	var book entities.Book
	var authorsUUID string
	var authors string

	for rows.Next() {
		err := rows.Scan(
			&book.ID,
			&book.UUID,
			&book.ISBN,
			&book.Title,
			&book.PublicationYear,
			&book.Description,
			&authors,
			&authorsUUID,
		)
		if err != nil {
			return nil, err
		}
		var a entities.Authors
		aList := strings.Split(authors, ", ")
		aUUID := strings.Split(authorsUUID, ", ")
		for authorNumber, author := range aList {
			a.Authors = append(a.Authors, entities.Author{
				UUID:     aUUID[authorNumber],
				NickName: author,
			})
		}
		book.BookAuthors.Authors = a.Authors
		books.Books = append(books.Books, book)
	}
	return &books, err
}

//func (r *bookRepository) gen() {
//	books := []struct {
//		NickName    string
//		Name        string
//		Surname     string
//		Patronymic  string
//		Description string
//		Birthday    string
//	}{
//
//		{"george_orwell", "George", "Orwell", "", "English novelist and essayist, journalist and critic.", "1903-06-25"},
//		{"harper_lee", "Harper", "Lee", "", "Author of 'To Kill a Mockingbird'.", "1926-04-28"},
//		{"f_scott_fitzgerald", "F. Scott", "Fitzgerald", "", "American novelist of the Jazz Age.", "1896-09-24"},
//		{"herman_melville", "Herman", "Melville", "", "American novelist, short story writer, and poet.", "1819-08-01"},
//		{"jane_austen", "Jane", "Austen", "", "English novelist known for her realism and biting social commentary.", "1775-12-16"},
//		{"j_d_salinger", "J.D.", "Salinger", "", "American writer known for 'The Catcher in the Rye'.", "1919-01-01"},
//		{"aldous_huxley", "Aldous", "Huxley", "", "English writer and philosopher.", "1894-07-26"},
//		{"jrr_tolkien", "J.R.R.", "Tolkien", "", "English writer, poet, and philologist, author of 'The Hobbit'.", "1892-01-03"},
//		{"leo_tolstoy", "Leo", "Tolstoy", "", "Russian writer, known for 'War and Peace'.", "1828-09-09"},
//		{"fyodor_dostoevsky", "Fyodor", "Dostoevsky", "", "Russian novelist and philosopher.", "1821-11-11"},
//		{"gabriel_garcia_marquez", "Gabriel", "Garcia Marquez", "", "Colombian novelist and Nobel laureate.", "1927-03-06"},
//		{"john_steinbeck", "John", "Steinbeck", "", "American author and Nobel Prize winner.", "1902-02-27"},
//		{"james_joyce", "James", "Joyce", "", "Irish novelist and poet.", "1882-02-02"},
//		{"miguel_de_cervantes", "Miguel", "de Cervantes", "", "Spanish writer, known for 'Don Quixote'.", "1547-09-29"},
//		{"homer", "", "Homer", "", "Ancient Greek poet, author of 'The Odyssey' and 'The Iliad'.", "-800-01-01"},
//		{"charles_dickens", "Charles", "Dickens", "", "English writer and social critic.", "1812-02-07"},
//		{"charlotte_bronte", "Charlotte", "Bronte", "", "English novelist and poet.", "1816-04-21"},
//		{"emily_bronte", "Emily", "Bronte", "", "English novelist and poet.", "1818-07-30"},
//		{"victor_hugo", "Victor", "Hugo", "", "French poet, novelist, and dramatist.", "1802-02-26"},
//		{"alexandre_dumas", "Alexandre", "Dumas", "", "French writer, known for historical novels.", "1802-07-24"},
//		{"gustave_flaubert", "Gustave", "Flaubert", "", "French novelist, author of 'Madame Bovary'.", "1821-12-12"},
//		{"dante_alighieri", "Dante", "Alighieri", "", "Italian poet, writer, and philosopher.", "1265-05-21"},
//		{"john_milton", "John", "Milton", "", "English poet and intellectual.", "1608-12-09"},
//		{"bram_stoker", "Bram", "Stoker", "", "Irish author, best known for 'Dracula'.", "1847-11-08"},
//		{"mary_shelley", "Mary", "Shelley", "", "English novelist, author of 'Frankenstein'.", "1797-08-30"},
//		{"oscar_wilde", "Oscar", "Wilde", "", "Irish poet and playwright.", "1854-10-16"},
//		{"nathaniel_hawthorne", "Nathaniel", "Hawthorne", "", "American novelist and short story writer.", "1804-07-04"},
//		{"mark_twain", "Mark", "Twain", "", "American writer and humorist.", "1835-11-30"},
//		{"ernest_hemingway", "Ernest", "Hemingway", "", "American novelist and journalist.", "1899-07-21"},
//		{"vladimir_nabokov", "Vladimir", "Nabokov", "", "Russian-American novelist.", "1899-04-22"},
//		{"toni_morrison", "Toni", "Morrison", "", "American novelist and Nobel laureate.", "1931-02-18"},
//		{"william_faulkner", "William", "Faulkner", "", "American writer and Nobel Prize winner.", "1897-09-25"},
//		{"ralph_ellison", "Ralph", "Ellison", "", "American novelist, author of 'Invisible Man'.", "1913-03-01"},
//		{"kurt_vonnegut", "Kurt", "Vonnegut", "", "American writer, known for 'Slaughterhouse-Five'.", "1922-11-11"},
//		{"joseph_heller", "Joseph", "Heller", "", "American author of 'Catch-22'.", "1923-05-01"},
//		{"ray_bradbury", "Ray", "Bradbury", "", "American author and screenwriter.", "1920-08-22"},
//		{"cormac_mccarthy", "Cormac", "McCarthy", "", "American novelist and playwright.", "1933-07-20"},
//		{"paulo_coelho", "Paulo", "Coelho", "", "Brazilian lyricist and novelist.", "1947-08-24"},
//		{"yann_martel", "Yann", "Martel", "", "Canadian author, best known for 'Life of Pi'.", "1963-06-25"},
//		{"khaled_hosseini", "Khaled", "Hosseini", "", "Afghan-American novelist.", "1965-03-04"},
//	}
//
//	for _, b := range books {
//		uid, _ := uuid.NewUUID()
//		res, err := r.db.ExecContext(context.Background(), `insert into authors (uuid, nickname, name, surname, patronymic, description, birth_date) values($1, $2, $3, $4, $5, $6, $7)`, uid.String(), b.NickName, b.Name, b.Surname, b.Patronymic, b.Description, b.Birthday)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(res)
//	}
//}

//func (r *bookRepository) gen() {
//	books := []struct {
//		Title           string
//		ISBN            string
//		PublicationYear int
//		Description     string
//		Publisher       string
//	}{
//		{"1984", "9780451524935", 1949, "A dystopian novel by George Orwell.", "Secker & Warburg"},
//		{"To Kill a Mockingbird", "9780061120084", 1960, "A novel about racial injustice in the Deep South.", "J.B. Lippincott & Co."},
//		{"The Great Gatsby", "9780743273565", 1925, "A novel about the American dream.", "Charles Scribner's Sons"},
//		{"Moby Dick", "9780142437247", 1851, "A novel about the voyage of the whaling ship Pequod.", "Harper & Brothers"},
//		{"Pride and Prejudice", "9780199535569", 1813, "A romantic novel by Jane Austen.", "T. Egerton"},
//		{"The Catcher in the Rye", "9780316769488", 1951, "A novel about teenage rebellion and alienation.", "Little, Brown and Company"},
//		{"Brave New World", "9780060850524", 1932, "A dystopian novel by Aldous Huxley.", "Chatto & Windus"},
//		{"The Hobbit", "9780547928227", 1937, "A fantasy novel by J.R.R. Tolkien.", "George Allen & Unwin"},
//		{"War and Peace", "9780199232765", 1869, "A novel by Leo Tolstoy about Russian society during the Napoleonic era.", "The Russian Messenger"},
//		{"Crime and Punishment", "9780486454115", 1866, "A psychological novel by Fyodor Dostoevsky.", "The Russian Messenger"},
//		{"Anna Karenina", "9780143035008", 1877, "A novel about Russian aristocracy by Leo Tolstoy.", "The Russian Messenger"},
//		{"The Brothers Karamazov", "9780374528379", 1880, "A novel by Fyodor Dostoevsky about morality and faith.", "The Russian Messenger"},
//		{"One Hundred Years of Solitude", "9780060883287", 1967, "A novel by Gabriel Garcia Marquez.", "Harper & Row"},
//		{"The Grapes of Wrath", "9780143039433", 1939, "A novel by John Steinbeck about the Great Depression.", "The Viking Press"},
//		{"Ulysses", "9780199535670", 1922, "A novel by James Joyce.", "Shakespeare and Company"},
//		{"Don Quixote", "9780060934347", 1605, "A novel by Miguel de Cervantes.", "Francisco de Robles"},
//		{"The Odyssey", "9780140268867", -800, "An epic poem attributed to Homer.", "Unknown"},
//		{"Great Expectations", "9780141439563", 1861, "A novel by Charles Dickens about personal growth.", "Chapman & Hall"},
//		{"Jane Eyre", "9780141441146", 1847, "A novel by Charlotte Bronte.", "Smith, Elder & Co."},
//		{"Wuthering Heights", "9780141439556", 1847, "A novel by Emily Bronte.", "Thomas Cautley Newby"},
//		{"Les Misérables", "9780451419439", 1862, "A novel by Victor Hugo.", "A. Lacroix, Verboeckhoven & Cie"},
//		{"The Count of Monte Cristo", "9780140449266", 1844, "A novel by Alexandre Dumas.", "Le Siècle"},
//		{"Madame Bovary", "9780140449129", 1857, "A novel by Gustave Flaubert.", "Revue de Paris"},
//		{"The Divine Comedy", "9780142437223", 1320, "An epic poem by Dante Alighieri.", "Various"},
//		{"Paradise Lost", "9780140424393", 1667, "An epic poem by John Milton.", "Samuel Simmons"},
//		{"The Iliad", "9780140275360", -750, "An epic poem attributed to Homer.", "Unknown"},
//		{"Dracula", "9780141439846", 1897, "A gothic novel by Bram Stoker.", "Archibald Constable & Co."},
//		{"Frankenstein", "9780141439471", 1818, "A novel by Mary Shelley about a scientist and his creation.", "Lackington, Hughes, Harding, Mavor & Jones"},
//		{"The Picture of Dorian Gray", "9780141439570", 1890, "A novel by Oscar Wilde about vanity and corruption.", "Lippincott's Monthly Magazine"},
//		{"A Tale of Two Cities", "9780141439600", 1859, "A novel by Charles Dickens set during the French Revolution.", "Chapman & Hall"},
//		{"The Scarlet Letter", "9780142437261", 1850, "A novel by Nathaniel Hawthorne.", "Ticknor, Reed & Fields"},
//		{"The Adventures of Huckleberry Finn", "9780143107323", 1884, "A novel by Mark Twain.", "Chatto & Windus"},
//		{"The Old Man and the Sea", "9780684801223", 1952, "A novel by Ernest Hemingway.", "Charles Scribner's Sons"},
//		{"Lolita", "9780679723165", 1955, "A controversial novel by Vladimir Nabokov.", "Olympia Press"},
//		{"Beloved", "9781400033416", 1987, "A novel by Toni Morrison about slavery and its aftermath.", "Alfred A. Knopf"},
//		{"The Sound and the Fury", "9780679732242", 1929, "A novel by William Faulkner.", "Jonathan Cape and Harrison Smith"},
//		{"Invisible Man", "9780679732761", 1952, "A novel by Ralph Ellison.", "Random House"},
//		{"Slaughterhouse-Five", "9780440180296", 1969, "A satirical novel by Kurt Vonnegut.", "Delacorte Press"},
//		{"Catch-22", "9781451626650", 1961, "A satirical novel by Joseph Heller.", "Simon & Schuster"},
//		{"Of Mice and Men", "9780140177398", 1937, "A novel by John Steinbeck.", "Covici Friede"},
//		{"Fahrenheit 451", "9781451673319", 1953, "A dystopian novel by Ray Bradbury.", "Ballantine Books"},
//		{"The Road", "9780307387899", 2006, "A post-apocalyptic novel by Cormac McCarthy.", "Alfred A. Knopf"},
//		{"The Alchemist", "9780062315007", 1988, "A philosophical novel by Paulo Coelho.", "HarperTorch"},
//		{"Life of Pi", "9780156027328", 2001, "A novel by Yann Martel.", "Knopf Canada"},
//		{"The Kite Runner", "9781594631931", 2003, "A novel by Khaled Hosseini.", "Riverhead Books"},
//	}
//	for _, b := range books {
//		uid, _ := uuid.NewUUID()
//		res, err := r.db.ExecContext(context.Background(), `insert into books (uuid, title, isbn, publication_year, description, publisher) values($1, $2, $3, $4, $5, $6)`, uid.String(), b.Title, b.ISBN, b.PublicationYear, b.Description, b.Publisher)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(res)
//	}
//}
