# Loan Engine Service

## Summary

This is a sample loan engine service that focus on RESTful API functionality with limited scope to create new loan, approve new loan, match investment (capital) with loan and disbursement. Other functionalities such as user registration, login, CRUD outside of loan objects will be outside the scope of this service.

### Within Scope

* create new loan
* approve loan
* create new investment to loan
* disburse loan

### Outside Scope and Assumptions

* no creation function for new user / employee / investor account
  * Assume to be handled by other service and there's existing accounts created for user, employee and investor
* no login and auth mechanism
  * Assume authentication is handled by other service and Loan service will use `X-User-Id` in API header to indicate the user who invoked the API
* user cannot submit arbitrary loan rate
  * Assume there's already fixed loan products with predefined rate and ROI.
* a loan can have multiple investments but an investment can only have one loan for simplicity
* loan cannot be editted
  * Assume loan amount, rate and roi are fixed after creation. The only attributes that get updated are the ones related to approval, invested and disbursed status update.
* investment cannot be editted
  * Assume once created investment cannot be changed
* the service does not handle agreement letter generation and storage
  * Assume it is created by other service and replaced by mocked functions that return dummy file storage URL
* there's no file upload for signed agreement letter
  * Assume it is already handled by other service and the API only receive the file URL of the agreement letter
* there's no verification and fraud checking
  * Assume once employee approve loan with picture proof, the loan will automatically be approved.
  * Similarly to disbursement, once employee upload signed agreement, it will automatically be disbursed.
* OpenTelemetry not implemented to simplify service PoC

## Product Specifications

### Loan States

```
proposed -> approved -> invested -> disbursed
```

Loan state can only move forward. It cannot be rolled back.

#### Expected Loan Flow
* user submit loan request via API
* employee approve loan and submit photo proof URL via API
* investor can make investment to a loan via API
  * loan will change state to `invested` only if the total amount of investment equal to loan amount
* employee disburse the loan and submit signed agreement document URL via API

### API Blueprint

API blueprint can be found on [/docs/OpenAPI.yaml](/docs/OpenAPI.yaml).

### Documentations

Diagram documentations can be found on [/docs](/docs/). All diagram are written in D2 Documentation format. You can use [D2 Playground](https://play.d2lang.com/) to view the diagram.

More info on [D2Lang](https://d2lang.com/).
