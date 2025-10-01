package org.example.repository;

import org.example.model.Notification;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface
NotificationRepository extends MongoRepository<Notification, String> {
    

}
