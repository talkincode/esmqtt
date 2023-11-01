package mqttc

// func PublicCustomerUpdate(customerid interface{}) error {
// 	var customer models.EdgeCustomer
// 	err := app.GDB().Where("id=?", customerid).First(&customer).Error
// 	if err != nil {
// 		return err
// 	}
// 	message := NewMessage[models.EdgeCustomer](MqttCustomerUpdate, customer)
// 	payload, err := message.Encode()
// 	if err != nil {
// 		return err
// 	}
// 	go Publish(1, MqttCustomerUpdate, payload)
// 	return nil
// }
