package server

import (
	"fmt"
	"log"
	"time"
)

func (s *APIServer) HandleRecurringTransactions() {
	r, err := s.DB.RecurringTransactions.FindAll_Internal()
	if err != nil {
		log.Println("ERROR IN CRON JOB: ", err)
	}

	for _, t := range r {
		if t.NextExecution.Before(time.Now().UTC()) {
			transaction := t.Transaction
			transaction.Date = time.Now()
			res, err := s.DB.Transactions.Save(&transaction)
			if err != nil {
				fmt.Println("ERROR ON RECURRING TRANSACTION with ID: ", t.ID, err)
			} else {
				// t.NextExecution
				switch t.UnitOfMeasure {
				case "month":
					t.NextExecution = t.LastExecution.AddDate(0, t.Frequency, 0)
				case "year":
					t.NextExecution = t.LastExecution.AddDate(t.Frequency, 0, 0)
				case "week":
					t.NextExecution = t.LastExecution.AddDate(0, 0, 7)
				case "day":
					t.NextExecution = t.LastExecution.AddDate(0, 0, t.Frequency)
				}

				t.LastExecution = time.Now().UTC()

				_, err := s.DB.RecurringTransactions.Save(t)

				if err != nil {
					fmt.Println("ERROR SAVING: ", err)
					fmt.Println("rolling back")
					if err := s.DB.Transactions.Delete(res.ID); err != nil {
						fmt.Println("err rolling back", err)
					}
				} else {
					fmt.Println("Next Execution: ", t.NextExecution)
				}
			}
		}
	}

}
