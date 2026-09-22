package repository

import "gorm.io/gorm/clause"

// clauseLockingUpdate is SELECT ... FOR UPDATE, used inside transactions to
// serialize concurrent reservation/release operations on the same pet.
var clauseLockingUpdate = clause.Locking{Strength: "UPDATE"}
