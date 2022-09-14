package openapi

type Visitor interface {
	VisitDocument(node *Document) (Visitor, error)
	VisitCallbacks(node *Callbacks) (Visitor, error)
	VisitComponents(node *Components) (Visitor, error)
	VisitContact(node *Contact) (Visitor, error)
	VisitDiscriminator(node *Discriminator) (Visitor, error)
}
