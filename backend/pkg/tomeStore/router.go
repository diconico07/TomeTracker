package tomeStore

import (
	"strconv"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/gin-gonic/gin"
)

func DefineRoutes(router *gin.RouterGroup, seriesSvc SeriesService, bookSvc BookService) {
	router.GET("/series", func(c *gin.Context) {
		series, err := seriesSvc.ListSeries(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		output := make([]struct {
			ent.Series
			VolumesTotal int `json:"volumes_total"`
			VolumesOwned int `json:"volumes_owned"`
		}, 0, len(series))
		for _, currentSeries := range series {
			volumesTotal, volumesOwned, err := seriesSvc.CountBooksInSeries(c.Request.Context(), currentSeries.ID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			output = append(output, struct {
				ent.Series
				VolumesTotal int "json:\"volumes_total\""
				VolumesOwned int "json:\"volumes_owned\""
			}{
				Series:       *currentSeries,
				VolumesTotal: volumesTotal,
				VolumesOwned: volumesOwned,
			})
		}

		c.JSON(200, output)
	})

	router.POST("/series", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		series, err := seriesSvc.AddSeries(c.Request.Context(), req.Name, req.URL)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, series)
	})

	router.DELETE("/series/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = seriesSvc.DeleteSeries(c.Request.Context(), id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})

	router.GET("/series/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		series, err := seriesSvc.GetSeries(c.Request.Context(), id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, series)
	})

	router.PATCH("/series/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var req SeriesUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		updatedSeries, err := seriesSvc.UpdateSeries(c.Request.Context(), id, req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, updatedSeries)
	})

	router.GET("/books", func(c *gin.Context) {
		books, err := bookSvc.ListBooks(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, books)
	})

	router.POST("/books", func(c *gin.Context) {
		var req struct {
			ISBN       string `json:"isbn"`
			Title      string `json:"title"`
			URL        string `json:"url"`
			SeriesID   int    `json:"series_id"`
			TomeNumber int    `json:"tome_number"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		book, err := bookSvc.AddBook(c.Request.Context(), req.Title, req.URL, req.SeriesID, req.TomeNumber, req.ISBN)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, book)
	})

	router.DELETE("/books/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := bookSvc.DeleteBook(c.Request.Context(), id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})

	router.GET("/books/:id", func(c *gin.Context) {
		id := c.Param("id")
		book, err := bookSvc.GetBook(c.Request.Context(), id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, book)
	})

	router.POST("/books/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req ent.Book
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		req.ID = id

		book, err := bookSvc.UpsertBook(c.Request.Context(), &req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, book)
	})

	router.PATCH("/books/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req BookUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		updatedBook, err := bookSvc.UpdateBook(c.Request.Context(), id, req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, updatedBook)
	})

	router.GET("/planning", func(c *gin.Context) {
		books, err := bookSvc.ListPlannedBooks(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, books)
	})

	router.GET("/missing", func(c *gin.Context) {
		books, err := bookSvc.ListMissingBooks(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, books)
	})

}
