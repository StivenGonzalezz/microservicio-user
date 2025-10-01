package org.example.controllers;


import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.example.dto.CreateNotificationDTO;
import org.example.dto.MessageDTO;
import org.example.model.Notification;
import org.example.service.interfaces.NotificationService;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/api")
@RequiredArgsConstructor
public class NotificationController {

    private final NotificationService service;

    @GetMapping("/notifications/channels")
    public List getChannels() {
        List channels = service.getChannels();
        return channels;
    }

    @PostMapping("/notifications")
    public ResponseEntity<MessageDTO> createNotification(@Valid @RequestBody CreateNotificationDTO notificationDTO)throws Exception {
        Notification notification = service.createNotification(notificationDTO);
        return ResponseEntity.ok(new MessageDTO<>(false,"Notification create successfully this ID: "+notification.getId()));
    }

    @GetMapping("/notifications")
    public ResponseEntity<MessageDTO<Page<Notification>>> getAllNotifications(@RequestParam(defaultValue = "0") int page, @RequestParam(defaultValue = "10") int size){
        Page<Notification> notifications = service.getNotifications(page,size);
        return ResponseEntity.ok(new MessageDTO<>(false,notifications));
    }

    @GetMapping("/notifications/{id}")
    public Notification getNotificationById(@PathVariable String id) throws Exception {
        Notification notification = service.getNotificationFindById(id);
        return notification;
    }

    @PostMapping("/notifications/schedule/{id}")
    public ResponseEntity<MessageDTO> scheduleNotification(@PathVariable String id) throws Exception{
        String message = service.scheduleNotification(id);
        return ResponseEntity.ok(new MessageDTO<>(false,message));
    }
}
