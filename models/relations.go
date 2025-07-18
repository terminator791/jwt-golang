package models

import "gorm.io/gorm"

// SetupRelations mengatur semua relasi antar model setelah semua tabel dibuat
func SetupRelations(db *gorm.DB) error {
	// Setup relasi User
	err := db.SetupJoinTable(&User{}, "Cards", &Card{})
	if err != nil {
		return err
	}

	// Setup relasi Terminal
	err = db.SetupJoinTable(&Terminal{}, "Gates", &Gate{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Terminal{}, "OriginTransactions", &Transaction{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Terminal{}, "DestTransactions", &Transaction{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Terminal{}, "FromFareMatrices", &FareMatrix{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Terminal{}, "ToFareMatrices", &FareMatrix{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Terminal{}, "SyncLogs", &SyncLog{})
	if err != nil {
		return err
	}

	// Setup relasi Card
	err = db.SetupJoinTable(&Card{}, "Transactions", &Transaction{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Card{}, "CardBalanceLogs", &CardBalanceLog{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Card{}, "TopUps", &TopUp{})
	if err != nil {
		return err
	}

	// Setup relasi Gate
	err = db.SetupJoinTable(&Gate{}, "CheckinTransactions", &Transaction{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&Gate{}, "CheckoutTransactions", &Transaction{})
	if err != nil {
		return err
	}

	// Setup relasi Transaction
	err = db.SetupJoinTable(&Transaction{}, "CardBalanceLog", &CardBalanceLog{})
	if err != nil {
		return err
	}

	return nil
}

// CreateForeignKeys membuat foreign key constraints setelah semua tabel dibuat
func CreateForeignKeys(db *gorm.DB) error {
	// Helper function untuk memeriksa apakah constraint sudah ada
	constraintExists := func(constraintName string) bool {
		var count int64
		db.Raw(`
			SELECT COUNT(*) 
			FROM information_schema.table_constraints 
			WHERE constraint_name = ? AND table_schema = current_schema()
		`, constraintName).Scan(&count)
		return count > 0
	}

	// Foreign keys untuk Card
	if !constraintExists("fk_users_cards") {
		err := db.Exec(`
			ALTER TABLE cards 
			ADD CONSTRAINT fk_users_cards 
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL;
		`).Error
		if err != nil {
			return err
		}
	}

	// Foreign keys untuk Gate
	if !constraintExists("fk_terminals_gates") {
		err := db.Exec(`
			ALTER TABLE gates 
			ADD CONSTRAINT fk_terminals_gates 
			FOREIGN KEY (terminal_id) REFERENCES terminals(terminal_id) ON DELETE CASCADE;
		`).Error
		if err != nil {
			return err
		}
	}

	// Foreign keys untuk Transaction
	if !constraintExists("fk_cards_transactions") {
		err := db.Exec(`
			ALTER TABLE transactions 
			ADD CONSTRAINT fk_cards_transactions 
			FOREIGN KEY (card_id) REFERENCES cards(card_id) ON DELETE CASCADE;
		`).Error
		if err != nil {
			return err
		}
	}

	if !constraintExists("fk_terminals_origin_transactions") {
		err := db.Exec(`
			ALTER TABLE transactions 
			ADD CONSTRAINT fk_terminals_origin_transactions 
			FOREIGN KEY (origin_terminal_id) REFERENCES terminals(terminal_id) ON DELETE SET NULL;
		`).Error
		if err != nil {
			return err
		}
	}

	if !constraintExists("fk_terminals_dest_transactions") {
		err := db.Exec(`
			ALTER TABLE transactions 
			ADD CONSTRAINT fk_terminals_dest_transactions 
			FOREIGN KEY (destination_terminal_id) REFERENCES terminals(terminal_id) ON DELETE SET NULL;
		`).Error
		if err != nil {
			return err
		}
	}

	if !constraintExists("fk_gates_checkin_transactions") {
		err := db.Exec(`
			ALTER TABLE transactions 
			ADD CONSTRAINT fk_gates_checkin_transactions 
			FOREIGN KEY (checkin_gate_id) REFERENCES gates(gate_id) ON DELETE SET NULL;
		`).Error
		if err != nil {
			return err
		}
	}

	if !constraintExists("fk_gates_checkout_transactions") {
		err := db.Exec(`
			ALTER TABLE transactions 
			ADD CONSTRAINT fk_gates_checkout_transactions 
			FOREIGN KEY (checkout_gate_id) REFERENCES gates(gate_id) ON DELETE SET NULL;
		`).Error
		if err != nil {
			return err
		}
	}

	// Foreign keys untuk TopUp
	if !constraintExists("fk_cards_top_ups") {
		err := db.Exec(`
			ALTER TABLE top_ups 
			ADD CONSTRAINT fk_cards_top_ups 
			FOREIGN KEY (card_id) REFERENCES cards(card_id) ON DELETE CASCADE;
		`).Error
		if err != nil {
			return err
		}
	}

	// Foreign keys untuk CardBalanceLog
	if !constraintExists("fk_cards_card_balance_logs") {
		err := db.Exec(`
			ALTER TABLE card_balance_logs 
			ADD CONSTRAINT fk_cards_card_balance_logs 
			FOREIGN KEY (card_id) REFERENCES cards(card_id) ON DELETE CASCADE;
		`).Error
		if err != nil {
			return err
		}
	}

	if !constraintExists("fk_transactions_card_balance_logs") {
		err := db.Exec(`
			ALTER TABLE card_balance_logs 
			ADD CONSTRAINT fk_transactions_card_balance_logs 
			FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id) ON DELETE SET NULL;
		`).Error
		if err != nil {
			return err
		}
	}

	// Foreign keys untuk SyncLog
	if !constraintExists("fk_terminals_sync_logs") {
		err := db.Exec(`
			ALTER TABLE sync_logs 
			ADD CONSTRAINT fk_terminals_sync_logs 
			FOREIGN KEY (terminal_id) REFERENCES terminals(terminal_id) ON DELETE CASCADE;
		`).Error
		if err != nil {
			return err
		}
	}

	// Foreign keys untuk FareMatrix
	if !constraintExists("fk_terminals_from_fare_matrices") {
		err := db.Exec(`
			ALTER TABLE fare_matrices 
			ADD CONSTRAINT fk_terminals_from_fare_matrices 
			FOREIGN KEY (from_terminal_id) REFERENCES terminals(terminal_id) ON DELETE CASCADE;
		`).Error
		if err != nil {
			return err
		}
	}

	if !constraintExists("fk_terminals_to_fare_matrices") {
		err := db.Exec(`
			ALTER TABLE fare_matrices 
			ADD CONSTRAINT fk_terminals_to_fare_matrices 
			FOREIGN KEY (to_terminal_id) REFERENCES terminals(terminal_id) ON DELETE CASCADE;
		`).Error
		if err != nil {
			return err
		}
	}

	return nil
}
