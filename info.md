# Changes to the files

## main.go
-   added to line 77
-   // Handler added to Read the Static files within the assets folder: .css, .js, img
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))

## html files
-   commented out the inline styles
-   added a script src link to the bottom of each file

## css
-   added a style.css file with a path of assests -> css -> style.css

## js
-   added a main.js file with a path of assests -> js -> main.js
-   added a script src link to the bottom of each html file