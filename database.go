package husky

import "gorm.io/gorm"

type DatabaseTxCall func(tx *gorm.DB) error

type DatabaseTxCalls interface {
	AddCall(call DatabaseTxCall)
	Run(db *gorm.DB) error
}

type _DatabaseTxCalls []DatabaseTxCall

func (calls *_DatabaseTxCalls) AddCall(call DatabaseTxCall) {
	*calls = append(*calls, call)
}

func (calls *_DatabaseTxCalls) Run(db *gorm.DB) error {
	if len(*calls) == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, call := range *calls {
			if err := call(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func NewDatabaseTxCalls() DatabaseTxCalls {
	return new(_DatabaseTxCalls)
}
