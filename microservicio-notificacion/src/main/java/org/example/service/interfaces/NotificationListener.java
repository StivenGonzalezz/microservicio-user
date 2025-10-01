package org.example.service.interfaces;

import org.example.dto.CreateNotificationDTO;

public interface NotificationListener {

    void receiveNotification (CreateNotificationDTO notificationDTO)throws Exception;
}
