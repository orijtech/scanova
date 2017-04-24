package scanova

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/orijtech/otils"
)

type Client struct {
	sync.RWMutex
	apiKey string

	rt http.RoundTripper
}

type Level string

const (
	LevelL Level = "L"
	LevelM Level = "M"
	LevelH Level = "H"
	LevelQ Level = "Q"
)

type Size string

const (
	SmallSize  Size = "s"
	MediumSize Size = "m"
	LargeSize  Size = "l"
	XLSize     Size = "xl"
	XXLSize    Size = "xxl"
	XXXLSize   Size = "xxxl"
)

type Response struct {
	rc      io.ReadCloser
	headers http.Header
}

var _ io.ReadCloser = (*Response)(nil)

func (r *Response) Read(b []byte) (int, error) { return r.rc.Read(b) }

func (r *Response) Headers() http.Header { return r.headers }

func (r *Response) Close() error { return r.rc.Close() }

type Request struct {
	URL string `json:"url,omitempty"`

	Logo *Logo `json:"logo,omitempty"`

	Size Size `json:"size,omitempty"`

	ErrorCorrection Level `json:"error_correction,omitempty"`

	EyePattern        EyeShape `json:"eye_pattern,omitempty"`
	DataGradientStyle string   `json:"data_gradient_style,omitempty"`

	InnerEyeColor   string `json:"eye_color_inner,omitempty"`
	OuterEyeColor   string `json:"eye_color_inner,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`

	DataGradientStartColor string `json:"data_gradient_start_color,omitempty"`
}

type DataPattern uint

type Gradient uint

const (
	GradientNone Gradient = iota
	GradientHorizontal
	GradientVertical
	GradientDiagonal
	GradientRadial
)

var translateGradients = map[Gradient]string{
	GradientNone:       "None",
	GradientHorizontal: "Horizontal",
	GradientVertical:   "Vertical",
	GradientDiagonal:   "Diagonal",
	GradientRadial:     "Radial",
}

type Logo struct {
	URL string `json:"url"`

	// Logo size in percentage
	PercentSize float32 `json:"size,omitempty"`

	Angle float32 `json:"angle,omitempty"`

	// Whether the logo should be excavated or not.
	// If true, all the data points overlapping with
	// the logo are removed.
	Excavated bool `json:"excavated"`
}

type Poster struct {
	URL string `json:"url"`

	PercentOfLeft float32 `json:"left,omitempty"`

	// Ratio of QR code size and image size in percentage.
	// (QR code size/ Image size) * 100.
	Size float32 `json:"size"`

	// Shape of eye of QR code
	EyeShape EyeShape `json:"eyeshape,omitempty"`
}

type Shape string

type EyeShape string

const (
	RoundRect              EyeShape = "ROUND_RECT"
	RectRect               EyeShape = "RECT_RECT"
	RectCircle             EyeShape = "RECT_CIRC"
	RoundRectangularCircle EyeShape = "ROUNDRECT_CIRC"
	CircularCircle         EyeShape = "CIRC_CIRC"

	BRLeaf         EyeShape = "BR_LEAF"
	TRLeaf         EyeShape = "TR_LEAF"
	BLLeaf         EyeShape = "BL_LEAF"
	TLLeaf         EyeShape = "TL_LEAF"
	TLBRLeaf       EyeShape = "TRBR_LEAF"
	TLBLLeaf       EyeShape = "TRBL_LEAF"
	TLBRLeafCircle EyeShape = "TRBL_LEAF_CIRC"
	TRBLLeafDIAD   EyeShape = "TRBL_LEAF_DIAD"
	RectDIAD       EyeShape = "RECT_DIAD"
	UniLeaf        EyeShape = "UNI_LEAF"
	BloatRect      EyeShape = "BLOAT_RECT"
	WarpRect0      EyeShape = "WARP_RECT0"
	CurveRect      EyeShape = "CURVE_RECT"
	DistRect       EyeShape = "DIST_RECT"
	ZigZag         EyeShape = "ZIGZAG"
	WarpRect1      EyeShape = "WARP_RECT1"
	BlackHole      EyeShape = "BLACK_HOLE"
	Star           EyeShape = "STAR"
	Grid           EyeShape = "GRID"
	Scion          EyeShape = "SCION"
	Octagon        EyeShape = "OCTAGON"
	Flower         EyeShape = "FLOWER"
	Hut            EyeShape = "HUT"
	DarkHut        EyeShape = "DARK_HUT"
)

type Pattern string

const (
	RoundPattern    Pattern = "ROUND"
	SquarePattern   Pattern = "SQUARE"
	DiamondsPattern Pattern = "DIAMONDS"
	OvalPattern     Pattern = "OVAL"
)

var errUnimplemented = errors.New("unimplemented")

var baseURL = "https://api.scanova.io/v2"

func (c *Client) NewQR(req *Request) (*Response, error) {
	urlValues, err := otils.ToURLValues(req)
	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s/qrcode/url?%s", baseURL, urlValues.Encode())
	res, err := c.httpClient().Get(fullURL)
	if err != nil {
		return nil, err
	}

	if !statusOK(res.StatusCode) {
		errMsg := res.Status

		// slurp the body first
		slurp, _ := ioutil.ReadAll(res.Body)
		if len(slurp) > 0 {
			errMsg = fmt.Sprintf("%s", slurp)
		}
		return nil, errors.New(errMsg)
	}

	resp := &Response{rc: res.Body, headers: res.Header}
	return resp, nil
}

func (c *Client) SetHTTPRoundTripper(rt http.RoundTripper) {
	c.Lock()
	c.rt = rt
	c.Unlock()
}

func (c *Client) httpClient() *http.Client {
	c.RLock()
	rt := c.rt
	c.RUnlock()

	if rt == nil {
		rt = http.DefaultTransport
	}

	return &http.Client{Transport: rt}
}

func NewClient() (*Client, error) {
	c := new(Client)
	return c, nil
}

func statusOK(code int) bool { return code >= 200 && code <= 299 }
