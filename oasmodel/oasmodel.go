package oasmodel

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

//refIndexElement part of map used to resolve $ref directive
type refIndexElement struct {
	name   string
	schema *SchemaOrRef
}

/*OpenAPI From OAS
openapi 		string 	REQUIRED. 			This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses. The openapi field SHOULD be used by tooling specifications and clients to interpret the OpenAPI document. This is not related to the API info.version string.
info 			Info Object 	REQUIRED. 	Provides metadata about the API. The metadata MAY be used by tooling as required.
servers 		[Server Object] 			An array of Server Objects, which provide connectivity information to a target server. If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of /.
paths 			Paths Object 	REQUIRED. 	The available paths and operations for the API.
components 		Components Object 			An element to hold various schemas for the specification.
security 		[Security Requirement Object] 	A declaration of which security mechanisms can be used across the API. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. Individual operations can override this definition. To make security optional, an empty security requirement ({}) can be included in the array.
tags 			[Tag Object] 				A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
externalDocs 	External Documentation Object 	Additional external documentation.
*/
type OpenAPI struct {
	Openapi      string                     `yaml:"openapi"`
	Info         Info                       `yaml:"info"`
	Servers      []Server                   `yaml:"servers,omitempty"`
	Paths        map[string]PathItem        `yaml:"paths"`
	Components   Components                 `yaml:"components,omitempty"`
	Security     []SecurityReq              `yaml:"security,omitempty"`
	Tags         []Tag                      `yaml:"tags,omitempty"`
	ExternalDocs ExternalDocs               `yaml:"externalDocs,omitempty"`
	RefIndex     map[string]refIndexElement `yaml:"-"`
	XWsRPC       map[string]XwsRPCService   `yaml:"x-ws-rpc,omitempty"`
}

/*Tag Object from OAS
name 			string 	REQUIRED. 		The name of the tag.
description 	string 					A short description for the tag. CommonMark syntax MAY be used for rich text representation.
externalDocs 	External 				Documentation Object 	Additional external documentation for this tag.
*/
type Tag struct {
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description,omitempty"`
	ExternalDocs ExternalDocs `yaml:"externalDocs,omitempty"`
}

/*Components Object from OAS
schemas 		Map[string, Schema Object | Reference Object] 		An object to hold reusable Schema Objects.
responses 		Map[string, Response Object | Reference Object] 	An object to hold reusable Response Objects.
parameters 		Map[string, Parameter Object | Reference Object] 	An object to hold reusable Parameter Objects.
examples 		Map[string, Example Object | Reference Object] 		An object to hold reusable Example Objects.
requestBodies 	Map[string, Request Body Object | Reference Object] 	An object to hold reusable Request Body Objects.
headers 		Map[string, Header Object | Reference Object] 		An object to hold reusable Header Objects.
securitySchemes 	Map[string, Security Scheme Object | Reference Object] 	An object to hold reusable Security Scheme Objects.
links 			Map[string, Link Object | Reference Object] 		An object to hold reusable Link Objects.
callbacks 		Map[string, Callback Object | Reference Object] 	An object to hold reusable Callback Objects.
*/
type Components struct {
	Schemas         map[string]*SchemaOrRef         `yaml:"schemas,omitempty"`
	Responses       map[string]*ResponseOrRef       `yaml:"responses,omitempty"`
	Parameters      map[string]*ParameterOrRef      `yaml:"parameters,omitempty"`
	Examples        map[string]*ExampleOrRef        `yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBodyOrRef    `yaml:"requestBodies,omitempty"`
	Headers         map[string]*HeaderOrRef         `yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeOrRef `yaml:"securitySchemes,omitempty"`
	Links           map[string]*LinkOrRef           `yaml:"links,omitempty"`
	Callbacks       map[string]*CallbackOrRef       `yaml:"callbacks,omitempty"`
}

/*Info ... From OAS Specifications :
title 			string REQUIRED. 	The title of the API.
description 	string 				A short description of the API. CommonMark syntax MAY be used for rich text representation.
termsOfService 	string 				A URL to the Terms of Service for the API. MUST be in the format of a URL.
contact			Contact Object		The contact information for the exposed API.
license 		License Object		The license information for the exposed API.
version			string REQUIRED. 	The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
*/
type Info struct {
	Title          string   `yaml:"title"`
	Description    string   `yaml:"description,omitempty"`
	TermsOfService string   `yaml:"termsOfService,omitempty"`
	Contact        *Contact `yaml:"contact,omitempty"`
	License        *License `yaml:"license,omitempty"`
	Version        string   `yaml:"version"`
}

/*Contact from OAS Specifications :
name      string The identifying name of the contact person/organization.
url       string The URL pointing to the contact information. MUST be in the format of a URL.
email     string The email address of the contact person/organization. MUST be in the format of an email address.
*/
type Contact struct {
	Name  string `yaml:"name,omitempty"`
	URL   string `yaml:"url,omitempty"`
	Email string `yaml:"email,omitempty"`
}

/*License from OAS Specifications :
name      string REQUIRED. The license name used for the API.
url string A URL to the license used for the API. MUST be in the format of a URL.
*/
type License struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url,omitempty"`
}

/*Server from OAS Specifications:
url string REQUIRED. A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
description string An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
variables Map[string, Server Variable Object] A map between a variable name and its value. The value is used for substitution in the server's URL template.
*/
type Server struct {
	URL         string                    `yaml:"url,omitempty"`
	Description string                    `yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `yaml:"variables,omitempty"`
}

/*ServerVariable from OAS Specifications:
enum [string] An enumeration of string values to be used if the substitution options are from a limited set. The array SHOULD NOT be empty.
default string REQUIRED. The default value to use for substitution, which SHALL be sent if an alternate value is not supplied. Note this behavior is different than the Schema Object's treatment of default values, because in those cases parameter values are optional. If the enum is defined, the value SHOULD exist in the enum's values.
description string An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
*/
type ServerVariable struct {
	Enum        []string `yaml:"enum,omitempty"`
	Default     string   `yaml:"default"`
	Description string   `yaml:"decription,omitempty"`
}

/*PathItem Object from OAS Specifications
/{path} Path Item Object A relative path to an individual endpoint. The field name MUST begin with a forward slash (/).
The path is appended (no relative URL resolution) to the expanded URL from the Server Object's url
field in order to construct the full URL. Path templating is allowed. When matching URLs, concrete (non-templated)
paths would be matched before their templated counterparts.
Templated paths with the same hierarchy but different templated names MUST NOT exist as they are identical.
In case of ambiguous matching, it's up to the tooling to decide which one to use.

$ref 		string 				Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object. In case a Path Item Object field appears both in the defined object and the referenced object, the behavior is undefined.
summary 	string 				An optional, string summary, intended to apply to all operations in this path.
description string 				An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
get			Operation Object	A definition of a GET operation on this path.
put			Operation Object	A definition of a PUT operation on this path.
post		Operation Object	A definition of a POST operation on this path.
delete		Operation Object	A definition of a DELETE operation on this path.
options		Operation Object	A definition of a OPTIONS operation on this path.
head		Operation Object	A definition of a HEAD operation on this path.
patch		Operation Object	A definition of a PATCH operation on this path.
trace		Operation Object	A definition of a TRACE operation on this path.
servers		[Server Object]		An alternative server array to service all operations in this path.
parameters	[Parameter Object | Reference Object] A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
*/
type PathItem struct {
	Ref         string      `yaml:"$ref,omitempty"`
	Summary     string      `yaml:"summary,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Get         *Operation  `yaml:"get,omitempty"`
	Put         *Operation  `yaml:"put,omitempty"`
	Post        *Operation  `yaml:"post,omitempty"`
	Delete      *Operation  `yaml:"delete,omitempty"`
	Options     *Operation  `yaml:"options,omitempty"`
	Head        *Operation  `yaml:"head,omitempty"`
	Patch       *Operation  `yaml:"patch,omitempty"`
	Trace       *Operation  `yaml:"trace,omitempty"`
	Servers     []Server    `yaml:"servers,omitempty"`
	Parameters  []Parameter `yaml:"parameters,omitempty"`
}

/*Operation Object from OAS
tags			[string]								A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
summary 		string									A short summary of what the operation does.
description		string									A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
externalDocs	External Documentation Object			Additional external documentation for this operation.
operationId		string									Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is case-sensitive. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
parameters 		[Parameter Object | Reference Object]	A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
requestBody 	Request Body Object | Reference Object	The request body applicable for this operation. The requestBody is only supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague, requestBody SHALL be ignored by consumers.
responses		Responses Object	REQUIRED. 			The list of possible responses as they are returned from executing this operation.
callbacks		Map[string, Callback Object | Reference Object]	A map of possible out-of band callbacks related to the parent operation. The key is a unique identifier for the Callback Object. Each value in the map is a Callback Object that describes a request that may be initiated by the API provider and the expected responses.
deprecated		boolean									Declares this operation to be deprecated. Consumers SHOULD refrain from usage of the declared operation. Default value is false.
security		[Security Requirement Object]			A declaration of which security mechanisms can be used for this operation. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. To make security optional, an empty security requirement ({}) can be included in the array. This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.
servers			[Server Object]							An alternative server array to service this operation. If an alternative server object is specified at the Path Item Object or Root level, it will be overridden by this value.
*/
type Operation struct {
	Tags         []string            `yaml:"tags,omitempty"`
	Summary      string              `yaml:"summary,omitempty"`
	Description  string              `yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs       `yaml:"externalDocs,omitempty"`
	OperationID  string              `yaml:"operationId,omitempty"`
	Parameters   []Parameter         `yaml:"parameters,omitempty"`
	RequestBody  *RequestBody        `yaml:"requestBody,omitempty"`
	Responses    Responses           `yaml:"responses"`
	Callbacks    map[string]Callback `yaml:"callbacks,omitempty"`
	Deprecated   bool                `yaml:"deprecated,omitempty"`
	Security     []SecurityReq       `yaml:"security,omitempty"`
	Servers      []Server            `yaml:"servers,omitempty"`
	XWsRPC       string              `yaml:"x-ws-rpc,omitempty"` // x-ws-rpc Link to Service WebSocket Open Operation
}

/*ExternalDocs from OAS
description	string	A short description of the target documentation. CommonMark syntax MAY be used for rich text representation.
url	string	REQUIRED. The URL for the target documentation. Value MUST be in the format of a URL.
*/
type ExternalDocs struct {
	Description string `yaml:"description,omitempty"`
	URL         string `yaml:"url"`
}

/*Parameter Object from OAS
name			string REQUIRED. 		The name of the parameter. Parameter names are case sensitive.
in 				string REQUIRED. 		The location of the parameter. Possible values are "query", "header", "path" or "cookie".
description		string					A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
required		boolean					Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
deprecated		boolean					Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
allowEmptyValue	boolean					Sets the ability to pass empty-valued parameters. This is valid only for query parameters and allows sending a parameter with an empty value. Default value is false. If style is used, and if behavior is n/a (cannot be serialized), the value of allowEmptyValue SHALL be ignored. Use of this property is NOT RECOMMENDED, as it is likely to be removed in a later revision.

style			string					Describes how the parameter value will be serialized depending on the type of the parameter value. Default values (based on value of in): for query - form; for path - simple; for header - simple; for cookie - form.
explode			boolean					When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this property has no effect. When style is form, the default value is true. For all other styles, the default value is false.
allowReserved	boolean					Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding. This property only applies to parameters with an in value of query. The default value is false.
schema			Schema Object | Reference Object	The schema defining the type used for the parameter.
example			Any					Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema that contains an example, the example value SHALL override the example provided by the schema. To represent examples of media types that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
examples		Map[ string, Example Object | Reference Object] Examples of the parameter's potential value. Each example SHOULD contain a value in the correct format as specified in the parameter encoding. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema that contains an example, the examples value SHALL override the example provided by the schema.
content			Map[string, Media Type Object]					A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry.
*/
type Parameter struct {
	Name            string                  `yaml:"name"`
	IN              string                  `yaml:"in"`
	Description     string                  `yaml:"description,omitempty"`
	Required        bool                    `yaml:"required,omitempty"`
	Deprecated      bool                    `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                    `yaml:"allowEmptyValue,omitempty"`
	Style           string                  `yaml:"style,omitempty"`
	Explode         *bool                   `yaml:"explode,omitempty"` // default value depends on style
	AllowReserved   bool                    `yaml:"allowReserved,omitempty"`
	Schema          *SchemaOrRef            `yaml:"schema,omitempty"`
	Example         *ExampleValue           `yaml:"example,omitempty"`
	Examples        map[string]ExampleOrRef `yaml:"examples,omitempty"`
	Content         map[string]MediaType    `yaml:"content,omitempty"`
}

/*RequestBody Object from OAS:
description	string A brief description of the request body. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
content Map[string, Media Type Object] REQUIRED. The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
required boolean Determines if the request body is required in the request. Defaults to false.
*/
type RequestBody struct {
	Ref         string               `yaml:"$ref,omitempty"`
	Description string               `yaml:"description,omitempty"`
	Content     map[string]MediaType `yaml:"content"`
	Required    bool                 `yaml:"required"`
}

/*SecurityScheme from OAS
type 				string 	Any 	REQUIRED. The type of the security scheme. Valid values are "apiKey", "http", "oauth2", "openIdConnect".
description 		string 	Any 				A short description for security scheme. CommonMark syntax MAY be used for rich text representation.
name 				string 	apiKey 	REQUIRED. 	The name of the header, query or cookie parameter to be used.
in 					string 	apiKey 	REQUIRED. 	The location of the API key. Valid values are "query", "header" or "cookie".
scheme 				string 	http 	REQUIRED. 	The name of the HTTP Authorization scheme to be used in the Authorization header as defined in RFC7235. The values used SHOULD be registered in the IANA Authentication Scheme registry.
bearerFormat 		string 	http ("bearer") 	A hint to the client to identify how the bearer token is formatted. Bearer tokens are usually generated by an authorization server, so this information is primarily for documentation purposes.
flows 				OAuth Flows Object 	oauth2 	REQUIRED. An object containing configuration information for the flow types supported.
openIdConnectUrl 	string 	openIdConnect 	REQUIRED. OpenId Connect URL to discover OAuth2 configuration values. This MUST be in the form of a URL.
*/
type SecurityScheme struct {
	Type             string     `yaml:"type"`
	Description      string     `yaml:"description,omitempty"`
	Name             string     `yaml:"name"`
	IN               string     `yaml:"in"`
	Scheme           string     `yaml:"scheme"`
	BearerFormat     string     `yaml:"bearerFormat,omitempty"`
	Flows            OAuthFlows `yaml:"flows"`
	OpenIDConnectURL string     `yaml:"openIdConnectUrl"`
}

/*OAuthFlows from OAS
implicit 			OAuth Flow Object 	Configuration for the OAuth Implicit flow
password 			OAuth Flow Object 	Configuration for the OAuth Resource Owner Password flow
clientCredentials 	OAuth Flow Object 	Configuration for the OAuth Client Credentials flow. Previously called application in OpenAPI 2.0.
authorizationCode 	OAuth Flow Object 	Configuration for the OAuth Authorization Code flow. Previously called accessCode in OpenAPI 2.0.
*/
type OAuthFlows struct {
	Implicit          *OAuth `yaml:"implicit,omitempty"`
	Password          *OAuth `yaml:"password,omitempty"`
	ClientCredentials *OAuth `yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuth `yaml:"authorizationCode,omitempty"`
}

/*OAuth from OAS
authorizationUrl 	string 					oauth2 ("implicit", "authorizationCode") 	REQUIRED. The authorization URL to be used for this flow. This MUST be in the form of a URL.
tokenUrl 			string 					oauth2 ("password", "clientCredentials", "authorizationCode") 	REQUIRED. The token URL to be used for this flow. This MUST be in the form of a URL.
refreshUrl 			string 					oauth2 	The URL to be used for obtaining refresh tokens. This MUST be in the form of a URL.
scopes 				Map[string, string] 	oauth2 	REQUIRED. The available scopes for the OAuth2 security scheme. A map between the scope name and a short description for it. The map MAY be empty.
*/
type OAuth struct {
	AuthorizationURL string            `yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `yaml:"scopes"`
}

/*MediaType Object from OAS:
schema			Schema Object | Reference Object	The schema defining the content of the request, response, or parameter.
example			Any									Example of the media type. The example object SHOULD be in the correct format as specified by the media type. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema which contains an example, the example value SHALL override the example provided by the schema.
examples		Map[ string, Example Object | Reference Object]	Examples of the media type. Each example object SHOULD match the media type and specified schema if present. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema which contains an example, the examples value SHALL override the example provided by the schema.
encoding		Map[string, Encoding Object]	A map between a property name and its encoding information. The key, being the property name, MUST exist in the schema as a property. The encoding object SHALL only apply to requestBody objects when the media type is multipart or application/x-www-form-urlencoded.
*/
type MediaType struct {
	Schema   *SchemaOrRef            `yaml:"schema,omitempty"`
	Example  ExampleValue            `yaml:"example,omitempty"`
	Examples map[string]ExampleOrRef `yaml:"examples,omitempty"`
	Encoding map[string]Encoding     `yaml:"encoding,omitempty"`
}

/*Responses from OAS
default 			Response Object | Reference Object 	The documentation of responses other than the ones declared for specific HTTP response codes. Use this field to cover undeclared responses. A Reference Object can link to a response that the OpenAPI Object's components/responses section defines.
Patterned Fields
Field Pattern 	Type 	Description
HTTP Status Code 	Response Object | Reference Object 	Any HTTP status code can be used as the property name, but only one property per code, to describe the expected response for that HTTP status code. A Reference Object can link to a response that is defined in the OpenAPI Object's components/responses section. This field MUST be enclosed in quotation marks (for example, "200") for compatibility between JSON and YAML. To define a range of response codes, this field MAY contain the uppercase wildcard character X. For example, 2XX represents all response codes between [200-299]. Only the following range definitions are allowed: 1XX, 2XX, 3XX, 4XX, and 5XX. If a response is defined using an explicit code, the explicit code definition takes precedence over the range definition for that code.
*/
type Responses map[string]*ResponseOrRef

/*Response from OAS
description		string REQUIRED. 								A short description of the response. CommonMark syntax MAY be used for rich text representation.
headers			Map[string, Header Object | Reference Object] 	Maps a header name to its definition. RFC7230 states header names are case insensitive. If a response header is defined with the name "Content-Type", it SHALL be ignored.
content 		Map[string, Media Type Object]					A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it. For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
links 			Map[string, Link Object | Reference Object]		A map of operations links that can be followed from the response. The key of the map is a short name for the link, following the naming constraints of the names for Component Objects.
*/
type Response struct {
	Description string                     `yaml:"description"`
	Headers     map[string]*HeaderOrRef    `yaml:"headers,omitempty"`
	Content     map[string]*MediaTypeOrRef `yaml:"content,omitempty"`
	Links       map[string]*LinkOrRef      `yaml:"links,omitempty"`
}

type Callback map[string]PathItem

//UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (e *CallbackOrRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*e = CallbackOrRef{}
	ref := Ref{}
	err := unmarshal(&ref)
	if err != nil || ref.Ref == "" {
		val := make(Callback)
		err = unmarshal(&val)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error un marshalling CallbackOrRef")
			return err
		}
		e.Val = &val
		return nil
	}
	e.Ref = &ref

	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *CallbackOrRef) MarshalYAML() (interface{}, error) {
	if e.Ref != nil {
		return e.Ref, nil
	}
	return e.Val, nil
}

type SecurityReq map[string][]string

type Header struct {
	Ref             string               `yaml:"$ref,omitempty"`
	Description     string               `yaml:"description,omitempty"`
	Required        bool                 `yaml:"required,omitempty"`
	Deprecated      bool                 `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                 `yaml:"allowEmptyValue,omitempty"`
	Style           string               `yaml:"style,omitempty"`
	Explode         *bool                `yaml:"required,omitempty"` // default value depends on style
	AllowReserved   bool                 `yaml:"allowReserved,omitempty"`
	Schema          *Schema              `yaml:"schema,omitempty"`
	Example         *ExampleValue        `yaml:"example,omitempty"`
	Examples        map[string]Example   `yaml:"examples,omitempty"`
	Content         map[string]MediaType `yaml:"content,omitempty"`
}

/*Link from OAS
operationRef		string							A relative or absolute URI reference to an OAS operation. This field is mutually exclusive of the operationId field, and MUST point to an Operation Object. Relative operationRef values MAY be used to locate an existing Operation Object in the OpenAPI definition.
operationId			string							The name of an existing, resolvable OAS operation, as defined with a unique operationId. This field is mutually exclusive of the operationRef field.
parameters			Map[string, Any | {expression}]	A map representing parameters to pass to an operation as specified with operationId or identified via operationRef. The key is the parameter name to be used, whereas the value can be a constant or an expression to be evaluated and passed to the linked operation. The parameter name can be qualified using the parameter location [{in}.]{name} for operations that use the same parameter name in different locations (e.g. path.id).
requestBody			Any | {expression}				A literal value or {expression} to use as a request body when calling the target operation.
description			string							A description of the link. CommonMark syntax MAY be used for rich text representation.
server				Server Object					A server object to be used by the target operation.
*/
type Link struct {
	Ref          string                 `yaml:"$ref,omitempty"`
	OperationRef string                 `yaml:"operationRef,omitempty"`
	OperationID  string                 `yaml:"operationId,omitempty"`
	Parameters   map[string]interface{} `yaml:"parameters,omitempty"`
	RequestBody  map[string]interface{} `yaml:"requestBody,omitempty"`
	Description  string                 `yaml:"description,omitempty"`
	Server       Server                 `yaml:"server,omitempty"`
}

/*Schema from OAS
--Properties as is from JSON Schema
title				string
multipleOf			integer
maximum
exclusiveMaximum
minimum
exclusiveMinimum
maxLength
minLength
pattern (This string SHOULD be a valid regular expression, according to the Ecma-262 Edition 5.1 regular expression dialect)
maxItems
minItems
uniqueItems
maxProperties
minProperties
required
enum
-- Properties with Adjustements from JSON Schema
type - 			Value MUST be a string. Multiple types via an array are not supported.
allOf - 		Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
oneOf - 		Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
anyOf - 		Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
not - 			Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
items - 		Value MUST be an object and not an array. Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema. items MUST be present if the type is array.
properties -	Property definitions MUST be a Schema Object and not a standard JSON Schema (inline or referenced).
additionalProperties - Value can be boolean or object. Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema. Consistent with JSON Schema, additionalProperties defaults to true.
description - 			CommonMark syntax MAY be used for rich text representation.
format - 				See Data Type Formats for further details. While relying on JSON Schema's defined formats, the OAS offers a few additional predefined formats.
default - 				The default value represents what would be assumed by the consumer of the input as the value of the schema if one is not provided. Unlike JSON Schema, the value MUST conform to the defined type for the Schema Object defined at the same level. For example, if type is string, then default can be "foo" but cannot be 1.

-- Additional Fields
nullable		boolean					A true value adds "null" to the allowed type specified by the type keyword, only if type is explicitly defined within the same Schema Object. Other Schema Object constraints retain their defined behavior, and therefore may disallow the use of null as a value. A false value leaves the specified or default type unmodified. The default value is false.
discriminator	Discriminator Object	Adds support for polymorphism. The discriminator is an object name that is used to differentiate between other schemas which may satisfy the payload description. See Composition and Inheritance for more details.
readOnly		boolean					Relevant only for Schema "properties" definitions. Declares the property as "read only". This means that it MAY be sent as part of a response but SHOULD NOT be sent as part of the request. If the property is marked as readOnly being true and is in the required list, the required will take effect on the response only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
writeOnly		boolean					Relevant only for Schema "properties" definitions. Declares the property as "write only". Therefore, it MAY be sent as part of a request but SHOULD NOT be sent as part of the response. If the property is marked as writeOnly being true and is in the required list, the required will take effect on the request only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
xml				XML Object				This MAY be used only on properties schemas. It has no effect on root schemas. Adds additional metadata to describe the XML representation of this property.
externalDocs	External Documentation Object	Additional external documentation for this schema.
example			Any						A free-form property to include an example of an instance for this schema. To represent examples that cannot be naturally represented in JSON or YAML, a string value can be used to contain the example with escaping where necessary.
deprecated		boolean					Specifies that a schema is deprecated and SHOULD be transitioned out of usage. Default value is false.

*/
type Schema struct {
	Type                 string                  `yaml:"type,omitempty"`
	Title                string                  `yaml:"title,omitempty"`
	MultipleOf           int                     `yaml:"multipleOf,omitempty"`
	Maximum              int                     `yaml:"maximum,omitempty"`
	ExclusiveMaximum     int                     `yaml:"exclusiveMaximum,omitempty"`
	Minimum              int                     `yaml:"minimum,omitempty"`
	ExclusiveMinimum     int                     `yaml:"exclusiveMinimum,omitempty"`
	MaxLength            int                     `yaml:"maxLength,omitempty"`
	MinLength            int                     `yaml:"minLength,omitempty"`
	Pattern              string                  `yaml:"pattern,omitempty"`
	MaxItems             int                     `yaml:"maxItems,omitempty"`
	MinItems             int                     `yaml:"minItems,omitempty"`
	UniqueItems          bool                    `yaml:"uniqueItems,omitempty"`
	MaxProperties        int                     `yaml:"maxProperties,omitempty"`
	MinProperties        int                     `yaml:"minProperties,omitempty"`
	Required             []string                `yaml:"required,omitempty"`
	Enum                 []string                `yaml:"enum,omitempty"`
	AllOf                []*SchemaOrRef          `yaml:"allOf,omitempty"`
	OneOf                []*SchemaOrRef          `yaml:"oneOf,omitempty"`
	AnyOf                []*SchemaOrRef          `yaml:"anyOf,omitempty"`
	Items                *SchemaOrRef            `yaml:"items,omitempty"`
	XPropertiesOrder     []string                `yaml:"x-properties-order"`
	Properties           map[string]*SchemaOrRef `yaml:"properties,omitempty"`
	AdditionalProperties *AdditionalProperties   `yaml:"additionalProperties,omitempty"`
	Description          string                  `yaml:"description,omitempty"`
	Format               string                  `yaml:"format,omitempty"`
	Default              string                  `yaml:"default,omitempty"`
	Nullable             bool                    `yaml:"nullable,omitempty"`
	Discriminator        *Discriminator          `yaml:"discriminator,omitempty"`
	ReadOnly             bool                    `yaml:"readonly,omitempty"`
	WriteOnly            bool                    `yaml:"writeOnly,omitempty"`
	XML                  XML                     `yaml:"xml,omitempty"`
	ExternalDocs         *ExternalDocs           `yaml:"externalDocs,omitempty"`
	Example              *Example                `yaml:"example,omitempty"`
	Deprecated           bool                    `yaml:"depreacated,omitempty"`
}

type AdditionalProperties struct {
	IsBool       bool
	BooleanValue bool
	Schema       *SchemaOrRef
}

/*Encoding from OAS
contentType			string			The Content-Type for encoding a specific property. Default value depends on the property type: for string with format being binary – application/octet-stream; for other primitive types – text/plain; for object - application/json; for array – the default is defined based on the inner type. The value can be a specific media type (e.g. application/json), a wildcard media type (e.g. image/*), or a comma-separated list of the two types.
headers				Map[string, Header Object | Reference Object]	A map allowing additional information to be provided as headers, for example Content-Disposition. Content-Type is described separately and SHALL be ignored in this section. This property SHALL be ignored if the request body media type is not a multipart.
style				string				Describes how a specific property value will be serialized depending on its type. See Parameter Object for details on the style property. The behavior follows the same values as query parameters, including default values. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
explode				boolean				When this is true, property values of type array or object generate separate parameters for each value of the array, or key-value-pair of the map. For other types of properties this property has no effect. When style is form, the default value is true. For all other styles, the default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
allowReserved		boolean				Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding. The default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
*/
type Encoding struct {
	ContentType   string                 `yaml:"contentType,omitempty"`
	Headers       map[string]HeaderOrRef `yaml:"headers,omitempty"`
	Style         string                 `yaml:"style,omitempty"`
	Explode       *bool                  `yaml:"required,omitempty"` // default value depends on style
	AllowReserved bool                   `yaml:"allowReserved,omitempty"`
}

/*Example from OAS
summary		string	Short description for the example.
description	string	Long description for the example. CommonMark syntax MAY be used for rich text representation.
value		Any		Embedded literal example. The value field and externalValue field are mutually exclusive. To represent examples of media types that cannot naturally represented in JSON or YAML, use a string value to contain the example, escaping where necessary.
externalValue	string	A URL that points to the literal example. This provides the capability to reference examples that cannot easily be included in JSON or YAML documents. The value field and externalValue field are mutually exclusive.
*/
type Example struct {
	Ref           string       `yaml:"$ref,omitempty"`
	Summary       string       `yaml:"summary,omitempty"`
	Description   string       `yaml:"description,omitempty"`
	Value         ExampleValue `yaml:"value,omitempty"`
	ExternalValue string       `yaml:"externalValue,omitempty"`
}

type ExampleValue map[string]interface {
}

/*Discriminator from OAS
propertyName	string	REQUIRED. The name of the property in the payload that will hold the discriminator value.
mapping	Map[string, string]	An object to hold mappings between payload values and schema names or references.
*/
type Discriminator struct {
	PropertyName string            `yaml:"propertyName"`
	Mapping      map[string]string `yaml:"mapping,omitempty"`
}

/*XML Object from OAS
name	string	Replaces the name of the element/attribute used for the described schema property. When defined within items, it will affect the name of the individual XML elements within the list. When defined alongside type being array (outside the items), it will affect the wrapping element and only if wrapped is true. If wrapped is false, it will be ignored.
namespace	string	The URI of the namespace definition. Value MUST be in the form of an absolute URI.
prefix	string	The prefix to be used for the name.
attribute	boolean	Declares whether the property definition translates to an attribute instead of an element. Default value is false.
wrapped	boolean	MAY be used only for an array definition. Signifies whether the array is wrapped (for example, <books><book/><book/></books>) or unwrapped (<book/><book/>). Default value is false. The definition takes effect only when defined alongside type being array (outside the items).
*/
type XML struct {
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
	Prefix    string `yaml:"prefix,omitempty"`
	Attribute bool   `yaml:"attribute,omitempty"`
	Wrapped   bool   `yaml:"wrapped,omitempty"`
}

type Ref struct {
	Ref      string      `yaml:"$ref,omitempty"`
	Resolved interface{} `yaml:"-"`
	RefName  string
}
type CallbackOrRef struct {
	Ref *Ref
	Val *Callback
}
type ExampleOrRef struct {
	Ref *Ref
	Val *Example
}
type HeaderOrRef struct {
	Ref *Ref
	Val *Header
}
type LinkOrRef struct {
	Ref *Ref
	Val *Link
}
type MediaTypeOrRef struct {
	Ref *Ref
	Val *MediaType
}
type ParameterOrRef struct {
	Ref *Ref
	Val *Parameter
}
type RequestBodyOrRef struct {
	Ref *Ref
	Val *RequestBody
}
type ResponseOrRef struct {
	Ref *Ref
	Val *Response
}
type SchemaOrRef struct {
	Ref *Ref
	Val *Schema
}
type SecuritySchemeOrRef struct {
	Ref *Ref
	Val *SecurityScheme
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *MediaTypeOrRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*e = MediaTypeOrRef{nil, nil}
	ref := Ref{}
	err := unmarshal(&ref)
	if err != nil || ref.Ref == "" {
		val := MediaType{}
		err = unmarshal(&val)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error un marshalling MediaTypeOrRef")
			return err
		}

		e.Val = &val
		return nil
	}
	e.Ref = &ref
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *MediaTypeOrRef) MarshalYAML() (interface{}, error) {
	if e.Ref != nil {
		return e.Ref, nil
	}
	return e.Val, nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *ParameterOrRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*e = ParameterOrRef{}
	ref := Ref{}
	err := unmarshal(&ref)
	if err != nil {
		val := Parameter{}
		err = unmarshal(&val)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error un marshalling CallbackOrRef")
			return err
		}
		*e.Val = val
	}
	*e.Ref = ref

	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *ParameterOrRef) MarshalYAML() (interface{}, error) {
	if e.Ref != nil {
		return e.Ref, nil
	}
	return e.Val, nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *ResponseOrRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*e = ResponseOrRef{nil, nil}
	ref := Ref{}
	err := unmarshal(&ref)
	if err != nil || ref.Ref == "" {
		val := Response{}
		err = unmarshal(&val)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error un marshalling ResponseOrRef")
			return err
		}

		e.Val = &val
		return nil
	}
	e.Ref = &ref
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *ResponseOrRef) MarshalYAML() (interface{}, error) {
	if e.Ref != nil {
		return e.Ref, nil
	}
	return e.Val, nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (s *SchemaOrRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*s = SchemaOrRef{nil, nil}
	ref := Ref{}
	err := unmarshal(&ref)
	if err != nil || ref.Ref == "" {
		val := Schema{}
		err = unmarshal(&val)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error un marshalling SchemaOrRef")
			return err
		}

		s.Val = &val
		return nil
	}
	s.Ref = &ref
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (s *SchemaOrRef) MarshalYAML() (interface{}, error) {
	if s.Ref != nil {
		return s.Ref, nil
	}
	return s.Val, nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *AdditionalProperties) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*e = AdditionalProperties{}
	e.IsBool = false
	boolValue := false
	err := unmarshal(&boolValue)
	if err != nil {
		schema := SchemaOrRef{}
		err = unmarshal(&schema)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error un marshalling schema")
			return err
		}
		e.Schema = &schema
	}
	e.BooleanValue = boolValue

	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (e *AdditionalProperties) MarshalYAML() (interface{}, error) {
	if e.IsBool {
		return e.BooleanValue, nil
	}
	return e.Schema, nil
}

//UnMarshal : build OpenAPI struct form buffer
func (oa *OpenAPI) UnMarshal(buffer []byte) (*OpenAPI, error) {
	err := yaml.Unmarshal(buffer, oa)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}
	return oa, nil
}

//Load Charge le fichier de spec d'interface
func (oa *OpenAPI) Load(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, oa)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	return nil
}

//Save dump la spec courante
func (oa *OpenAPI) Save() {
	buf, err := yaml.Marshal(oa)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	fmt.Fprintf(os.Stdout, "%s", string(buf))
}

// $ref resolver
func (oa *OpenAPI) makeSchemaRefIndex() map[string]refIndexElement {
	refIndex := make(map[string]refIndexElement)
	for k, v := range oa.Components.Schemas {
		refIndex["#/components/schemas/"+k] = refIndexElement{k, v}
		v.fillRefIndex("#/components/schemas/"+k, "types."+k, refIndex)
	}
	return refIndex
}

// ResolveRefs
func (oa *OpenAPI) ResolveRefs() {
	refIndex := oa.makeSchemaRefIndex()
	for _, v := range oa.Components.Schemas {
		v.resolveRefs(refIndex)
	}
}

func (s *SchemaOrRef) fillRefIndex(yPath string, path string, refIndex map[string]refIndexElement) {
	if s.Ref != nil {
		return
	}
	if s.Val.Properties != nil {
		for p, v := range s.Val.Properties {
			// add property to Index
			refIndex[yPath+"/Properties/"+p] = refIndexElement{path + "." + p, v}
			v.fillRefIndex(yPath+"/"+p, path+"."+p, refIndex)
		}
	}
	/* TODO Doit for
	AllOf                []*SchemaOrRef          `yaml:"allOf,omitempty"`
	OneOf                []*SchemaOrRef          `yaml:"oneOf,omitempty"`
	AnyOf                []*SchemaOrRef          `yaml:"anyOf,omitempty"`
	Items                *SchemaOrRef            `yaml:"items,omitempty"`
	*/

}
func (s *SchemaOrRef) resolveRefs(refIndex map[string]refIndexElement) {
	if s.Ref != nil {
		log.Printf("Resolving %s ...\n", s.Ref.Ref)
		if elem, ok := refIndex[s.Ref.Ref]; ok {
			log.Printf("Resolving %s int %v\n", s.Ref.Ref, elem.schema)
			s.Ref.Resolved = elem.schema
			s.Ref.RefName = elem.name
		} else {
			log.Printf("Can't Resolve %s\n", s.Ref.Ref)
		}
		return
	}

	if s.Val.Properties != nil {
		for p, v := range s.Val.Properties {
			log.Printf("visit %s ...\n", p)
			v.resolveRefs(refIndex)
		}
	}
	if s.Val.Items != nil {
		s.Val.Items.resolveRefs(refIndex)
	}
	if s.Val.AllOf != nil {
		for p, v := range s.Val.AllOf {
			log.Printf("visit %d %v ...\n", p, v)
			v.resolveRefs(refIndex)
		}
	}
	if s.Val.OneOf != nil {
		for p, v := range s.Val.OneOf {
			log.Printf("visit %d %v ...\n", p, v)
			v.resolveRefs(refIndex)
		}
	}
	if s.Val.AnyOf != nil {
		for p, v := range s.Val.OneOf {
			log.Printf("visit %d %v ...\n", p, v)
			v.resolveRefs(refIndex)
		}
	}
	if s.Val.AdditionalProperties != nil {
		s.Val.AdditionalProperties.Schema.resolveRefs(refIndex)
	}
}
func (s *SchemaOrRef) Schema() *Schema {
	if s.Ref != nil && s.Ref.Resolved != nil {
		return s.Ref.Resolved.(*SchemaOrRef).Schema()
	}
	if s.Val != nil {
		return s.Val
	}
	log.Fatalf("unable to convert %s into valid schema", s.Ref.Ref)
	return nil
}
