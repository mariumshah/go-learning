package converter

import (
	// "reflect"
	"testing"

	"github.com/playground/userapi/internal/app/domain"
	"github.com/playground/userapi/pkg/models"
)

/*

We test:
- That ToDomain maps fields from the sqlboiler model into the domain model correctly.
- That ToModel maps domain model back into a sqlboiler model correctly,
  and that ID handling behaves as expected for both "update" (ID set) and "create" (ID == 0) cases.
*/

func TestConversionRoundTrip_UpdateCase(t *testing.T) {
	m := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	println("just checking", m.Email)

	d := UserToDomain(m)
	println("to domain: ", d.Email)
	// Basic assertions: fields must match
	if d.ID != m.ID {
		t.Fatalf("ID mismatch after ToDomain: got %d want %d", d.ID, m.ID)
	}
	if d.Email != m.Email {
		t.Fatalf("Email mismatch after ToDomain: got %q want %q", d.Email, m.Email)
	}
	if d.Hash != m.PasswordHash {
		t.Fatalf("PasswordHash mismatch after ToDomain: got %q want %q", d.Hash, m.PasswordHash)
	}

	m2 := UserToModel(d)

	if m2.ID != m.ID {
		t.Fatalf("ID mismatch after ToModel: got %d want %d", m2.ID, m.ID)
	}
	if m2.Email != m.Email {
		t.Fatalf("Email mismatch after ToModel: got %q want %q", m2.Email, m.Email)
	}
	if m2.PasswordHash != m.PasswordHash {
		t.Fatalf("PasswordHash mismatch after ToModel: got %q want %q", m2.PasswordHash, m.PasswordHash)
	}
}

func TestConversionCreateCase_IDZero(t *testing.T) {
	// Create a domain.User representing a new user (ID == 0)
	newDomain := &domain.User{
		// ID:    1, // new user -> DB should assign ID
		Email: "new@example.com",
		Hash:  "bcrypt$newhash",
	}

	// Convert to model for insertion
	modelForInsert := UserToModel(newDomain)
	println("ID", modelForInsert.ID)

	// On create we expect the model's ID to remain zero-value so DB will auto-increment it.
	// We only check the ID here; other fields must still be mapped.
	if modelForInsert.ID != 0 {
		t.Fatalf("expected model ID to be zero for create-case, got %d", modelForInsert.ID)
	}
	if modelForInsert.Email != newDomain.Email {
		t.Fatalf("Email mismatch on create-case: got %q want %q", modelForInsert.Email, newDomain.Email)
	}
	if modelForInsert.PasswordHash != newDomain.Hash {
		t.Fatalf("PasswordHash mismatch on create-case: got %q want %q", modelForInsert.PasswordHash, newDomain.Hash)
	}
}
