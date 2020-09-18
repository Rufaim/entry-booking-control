# Entry booking control application

Here is an pplication to book something on weekday basis.
Initially it is desined to laboratory entry booking, but might be applied to control booking of any shared resource.

**Attention:** in order to use application both server and client are to be put on the same physical server that all your users have access to.
For instance, it could be a shared computational resource.

## Building from scratch
*(next steps are valid for Linux system only)*

 First, clone repository:
 ```bash
git clone https://github.com/Rufaim/entry-booking-control.git
cd entry-booking-control
 ```
 
 Second, run make file:
 ```bash
 make
 ```
 
 It'll build both client and server parts.
 
 ## Usage
 Do not forget to start server first
 
 ```bash
 ./entry_booking_server
 ```
 Capacity and port could be changed in [here](/cmd/server/constants.go)
 
 ### Booking 
 To book a day run:
```bash
./entry_booking_client book -weekday <DAY>
```
where `<DAY>` is one of `Mon`, `Tue`, `Wed`, `Thu`, `Fri`.

To remove booking on a day run 

```bash
./entry_booking_client book -d -weekday <DAY>
```

### Checking visits
To watch you bookings do:
```bash
./entry_booking_client visits
```
*Example output*
```
user:
Mon   Wed   Thu
```
To watch full report on bookings run:
```bash
./entry_booking_client visits -a
```
*Example output*
```
Mon:
  first_user
  third_user
Tue:
  second_user
  third_user
Wed:
  first_user
  second_user
  third_user
Thu:
  first_user

Booking info:
Mon:2   Tue:2   Wed:3   Thu:1   Fri:0
```
