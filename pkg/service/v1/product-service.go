package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/MartyKuentzel/projectX/pkg/api/v1"
	"github.com/MartyKuentzel/projectX/pkg/logger"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// productServiceServer is implementation of v1.ProductServiceServer proto interface
type productServiceServer struct {
	db *sql.DB
}

// NewProductServiceServer creates Product service
func NewProductServiceServer(db *sql.DB) v1.ProductServiceServer {
	return &productServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *productServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *productServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// initialize table Product
func (s *productServiceServer) createTable(ctx context.Context, c *sql.Conn) error {

	_, err := c.ExecContext(ctx, "CREATE TABLE `Product` (`ID` bigint(20) NOT NULL AUTO_INCREMENT,"+
		"`Name` varchar(200) DEFAULT NULL,"+
		"`Price` varchar(200) DEFAULT NULL,"+
		"`Creator` varchar(200) DEFAULT NULL,"+
		"`Unit` varchar(200) DEFAULT NULL,"+
		"`Category` varchar(200) DEFAULT NULL,"+
		"`Description` varchar(1024) DEFAULT NULL,"+
		"`Date` timestamp NULL DEFAULT NULL,"+
		"PRIMARY KEY (`ID`),"+
		"UNIQUE KEY `ID_UNIQUE` (`ID`))")

	if err != nil {
		return status.Error(codes.Unknown, "failed to create table -> "+err.Error())
	}
	return nil
}

// Create new product task
func (s *productServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	date, err := ptypes.Timestamp(req.Product.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "date field has invalid format-> "+err.Error())
	}
	_, err = c.ExecContext(ctx, "SELECT 1 FROM Product LIMIT 1 ;")

	if err != nil {
		logger.Log.Warn("Table 'Product' doesn't exist: It will be created now.")
		err = s.createTable(ctx, c)
		if err != nil {
			return nil, err
		}
	}

	// insert Product entity data
	res, err := c.ExecContext(ctx, "INSERT INTO Product(`Name`, `Price`, `Creator`, `Unit`, `Category`, `Description`, `Date`) VALUES(?, ?, ?, ?, ?, ?, ?)",
		req.Product.Name, req.Product.Price, req.Product.Creator, req.Product.Unit, req.Product.Category, req.Product.Description, date)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Product-> "+err.Error())
	}

	// get ID of creates Product
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created Product-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read product task
func (s *productServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query product by ID
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Name`, `Price`, `Creator`, `Unit`, `Category`, `Description`, `Date` FROM Product WHERE `ID`=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Product-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Product-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Product with ID='%d' is not found",
			req.Id))
	}

	// get Product data
	var td v1.ProductProto
	var date time.Time
	if err := rows.Scan(&td.Id, &td.Name, &td.Price, &td.Creator, &td.Unit, &td.Category, &td.Description, &date); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Product row-> "+err.Error())
	}
	td.Date, err = ptypes.TimestampProto(date)
	if err != nil {
		return nil, status.Error(codes.Unknown, "date field has invalid format-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Product rows with ID='%d'",
			req.Id))
	}

	return &v1.ReadResponse{
		Api:     apiVersion,
		Product: &td,
	}, nil

}

// Update product task
func (s *productServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	date, err := ptypes.Timestamp(req.Product.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "date field has invalid format-> "+err.Error())
	}

	// update Product
	res, err := c.ExecContext(ctx, "UPDATE Product SET `Name`=?, `Price`=?, `Unit`=?, `Category`=?, `Creator`=?, `Description`=?, `Date`=? WHERE `ID`=?",
		req.Product.Name, req.Product.Price, req.Product.Unit, req.Product.Category, req.Product.Creator, req.Product.Description, date, req.Product.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update Product-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Product with ID='%d' is not found",
			req.Product.Id))
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil
}

// Delete product task
func (s *productServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete Product
	res, err := c.ExecContext(ctx, "DELETE FROM Product WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete Product-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Product with ID='%d' is not found",
			req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}

// Read all Product tasks
func (s *productServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get Product list
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Name`, `Price`, `Creator`, `Unit`, `Category`, `Description`, `Date` FROM Product")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Product-> "+err.Error())
	}
	defer rows.Close()

	var date time.Time
	list := []*v1.ProductProto{}
	for rows.Next() {
		td := new(v1.ProductProto)
		if err := rows.Scan(&td.Id, &td.Name, &td.Price, &td.Unit, &td.Category, &td.Creator, &td.Description, &date); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from Product row-> "+err.Error())
		}
		td.Date, err = ptypes.TimestampProto(date)
		if err != nil {
			return nil, status.Error(codes.Unknown, "date field has invalid format-> "+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Product-> "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:      apiVersion,
		Products: list,
	}, nil
}
