package org.example.service.interfaces;

import org.example.dto.CreateNotificationDTO;
import org.example.model.Notification;
import org.example.model.enums.Channels;
import org.springframework.data.domain.Page;

import java.util.List;


public interface NotificationService {

    Page<Notification> getNotifications(int page, int size) ;

    Notification getNotificationFindById(String id) throws Exception;

    List<Channels> getChannels();

    Notification createNotification(CreateNotificationDTO createNotificationDTO) throws Exception;

    String scheduleNotification(String id) throws Exception;


}


